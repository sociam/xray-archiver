/* global angular, _, $, d3 */

angular.module('dci')
	.controller('sankey', function () {})
	.config(function ($stateProvider, $urlRouterProvider) {
			$stateProvider.state('sankey', {
				url: '/sankey?app',
				template:'<div id="sankey-chart"></div>',
				resolve: {
					pitypes:($http) => $http.get('../mitm_out/pi_by_host.json').then((x) => x.data),
					hosts: ($http) => $http.get('../mitm_out/host_by_app.json').then((x) => x.data),
					details: ($http) => $http.get('../mitm_out/company_details.json').then((x) => x.data),
					data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
				},
				controller:function($scope, pitypes, hosts, details, data, $stateParams) {
					$scope.apps = _.uniq(data.map((x) => x.app));
					data = $scope.data = data.filter((x) => x.app === $stateParams.app);

					var margin = {top: 1, right: 1, bottom: 6, left: 1},
						width = 960 - margin.left - margin.right,
						height = 500 - margin.top - margin.bottom;

					var formatNumber = d3.format(",.0f"),
						format = function(d) { return formatNumber(d) + " TWh"; },
						color = d3.scale.category20();

					var svg = d3.select("#sankey-chart").append("svg")
					    .attr("width", width + margin.left + margin.right)
					    .attr("height", height + margin.top + margin.bottom)
					  .append("g")
					    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

					var sankey = d3.sankey()
					    .nodeWidth(15)
					    .nodePadding(10)
					    .size([width, height]),
					    path = sankey.link(),
					    link,
						dragmove = function(d) {
						    d3.select(this).attr("transform", "translate(" + d.x + "," + (d.y = Math.max(0, Math.min(height - d.dy, d3.event.y))) + ")");
						    sankey.relayout();
						    link.attr("d", path);
					    };

					if (!data.length) { $scope.error = 'no data for app ' + $stateParams.app; }

					var app = $scope.app = $stateParams.app,
						appcompany = $scope.appcompany = data[0].company,
						hTc = $scope.hTc = data.reduce((r,a) => {
							if (a.host_company) { 
								r[a.host] = a.host_company; 
								r[a.host_2ld] = a.host_company;
							}
							return r;
						}, {}),
						hTh = $scope.hTh = data.reduce((r,a) => {
							// host -> 2ld
							if (a.host_2ld) { 
								r[a.host] = a.host_2ld;
							} else { console.error('warning no 2ld ', a.host); }
							return r;
						}, {}),
						matchCompany = (x) => appcompany && ((x || '').toLowerCase() === appcompany.toLowerCase()),
						isType = $scope.isType = (id, type) => id && 
							details[id] && 
							details[id].typetag && 
							details[id].typetag.indexOf(type) >= 0,
						getName = $scope.getName = (id) => details[id] && details[id].company || id,
						is3rdPartyType = $scope.is3rdPartyType = (id, type) => isType(id,type) &&
								!_.some([getName(id), id].map(matchCompany)),
						recompute = () => {
							var apphosts = _(hosts[$scope.app]).keys().value();
							// next we wanna group together all the pi_types, and consolidate around company
							// console.info('threshold', $scope.threshold, 'apphosts', apphosts.length);
							$scope.company2pi = apphosts.reduce((r,host) => {
								var company = hTc[host], 
									host_pis = pitypes[host] || [];
								if (!company) { 
									var mfirst = hTh[host].match(/^([^\.]+)\./);
									if (mfirst) { 
										company = mfirst[1]; 
									} else {
										console.error('no company for host ', host); return r; 
									}
								} 
								r[company] = _.union(r[company] || [], host_pis);
								return r;
							}, {});

							// each of the boxes
							$scope.categories = {
								'app-publisher': _.pickBy($scope.company2pi, (pis, company) => 
										matchCompany(company)),
								'app-functionality': _.pickBy($scope.company2pi, (pis, company) => 
										!isType(company, 'ignore') && 
										!isType(company, 'platform') && 
										is3rdPartyType(company,'app')),
								'marketing': _.pickBy($scope.company2pi, (pis, company) => 
										!isType(company, 'ignore') && 
										!isType(company, 'platform') && 
										is3rdPartyType(company,'marketing')),
								'usage tracking': _.pickBy($scope.company2pi, (pis, company) => 
										!isType(company, 'ignore') && 
										!isType(company, 'platform') && 
										is3rdPartyType(company, 'usage')),
								'payments':_.pickBy($scope.company2pi, (pis, company) => 
										!isType(company, 'ignore') && 
										!isType(company, 'platform') && 
										is3rdPartyType(company, 'payments')),
								'security':_.pickBy($scope.company2pi, (pis, company) => 
										!isType(company, 'ignore') && 
										!isType(company, 'platform') && 
										is3rdPartyType(company, 'security')),
								'other': _.pickBy($scope.company2pi, (pis, company) => 
									!matchCompany(company) &&							
									!isType(company, 'ignore') && 							
									!isType(company, 'app') && 							
									!isType(company,'marketing') &&
									!isType(company, 'platform') && 							
									!isType(company, 'usage') &&
									!isType(company, 'payments') &&
									!isType(company, 'security'))
							};

							// let's start making nodes
							// start with the pitypes
							var OTHERPITYPE = 'OTHER_PI',
								pitypes_set = _(pitypes).values().flatten().uniq().value().concat([OTHERPITYPE]),
								nodemap = $scope.nodemap = {},
								nodes = $scope.nodes = [], 
								links = $scope.links = [],
								pushNode = (id, label) => {
									var l = nodes.length;
									nodes.push({name:id, label:label});
									nodemap[id] = l;
									return l;
								},
								pushLink = (from, to, width) => {
									var l = links.length;
									links.push({source:from, target:to, value:width || 1});
									return l;
								},
								pilabels = {
									USER_PERSONAL_DETAILS: 'personal details',
									USER_LOCATION: 'your location',
									USER_COARSE_LOCATION: 'your city/town',
									DEVICE_ID:'phone id',
									DEVICE_SOFT:'phone characteristics'
								};

							pitypes_set.map((pitype) => pushNode(pitype, pilabels[pitype]));

							// make nodes for categories as well 
							_.keys($scope.categories).map((cname) => pushNode(cname));

							// next we want to link these types to companies
							_.toPairs($scope.company2pi).map((pair) => {
								var company = pair[0], 
									pitypes = pair[1],
									company_nid = pushNode(company);

								pitypes.map((pitype) => { 
									var pitype_id = nodemap[pitype];
									console.info('adding pitype-company link ', pitype, pitype_id, ' -> ', company, company_nid);
									pushLink(pitype_id, company_nid);
								});
								_.keys($scope.categories).filter((cname) => $scope.categories[cname][company]).map((cname) => {
									console.info('adding purpose link ', company, company_nid, ' -> ', cname, nodemap[cname]);
									pushLink(company_nid, nodemap[cname], $scope.company2pi[company].length || 1);
								});

								if ($scope.company2pi[company].length === 0) {
									console.info('adding other_PI link ', nodemap[OTHERPITYPE], ' -> ', company, ' ', company_nid);
									pushLink(nodemap[OTHERPITYPE], company_nid);
								}
							});
							console.log('old nodes ', nodes.length, JSON.stringify({nodes:nodes}));
							console.log('old links ', links.length, JSON.stringify({links:links}));							

							// kill nodes that don't have any links
							var newnodes = _.filter(nodes, (n,i) => links.filter((l) => l.source === i || l.target === i).length > 0),
								newmap = $scope.newmap = newnodes.reduce((m, i) => { 
									m[i.name] = newnodes.indexOf(i);
									return m;
								}, {}),
								newlinks = $scope.newlinks = links.map((oldlink) =>_.extend({}, oldlink, {
									source: newmap[nodes[oldlink.source].name],
									target: newmap[nodes[oldlink.target].name]
								}));
	
							console.log('new nodes ', newnodes.length, JSON.stringify({nodes:newnodes}));
							console.log('new links ', newlinks.length, JSON.stringify({links:newlinks}));

							// do it
							sankey.nodes(newnodes)
							    .links(newlinks)
							    .layout(32);

							link = svg.append("g").selectAll(".link")
							  .data(newlinks)
							  .enter().append("path")
							  .attr("class", "link")
							  .attr("d", path)
							  .style("stroke-width", function(d) { return Math.max(1, d.dy); })
							  .sort(function(a, b) { return b.dy - a.dy; });

							link.append("title")
						    	.text(function(d) { return d.source.name + " → " + d.target.name; });

							var node = svg.append("g").selectAll(".node")
								.data(newnodes)
								.enter().append("g")
									.attr("class", "node")
									.attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"; })
								.call(d3.behavior.drag()
							  		.origin(function(d) { return d; })
							  		.on("dragstart", function() { this.parentNode.appendChild(this); })
							  		.on("drag", dragmove));

							node.append("rect")
							      .attr("height", function(d) { return d.dy; })
							      .attr("width", sankey.nodeWidth())
							      .style("fill", function(d) { return d.color = color(d.name.replace(/ .*/, "")); })
							      .style("stroke", function(d) { return d3.rgb(d.color).darker(2); })
							    .append("title")
							      .text(function(d) { return d.name + "\n" + format(d.value); });

							  node.append("text")
							      .attr("x", -6)
							      .attr("y", function(d) { return d.dy / 2; })
							      .attr("dy", ".35em")
							      .attr("text-anchor", "end")
							      .attr("transform", null)
							      .text(function(d) { return d.label || d.name; })
							    .filter(function(d) { return d.x < width / 2; })
							      .attr("x", 6 + sankey.nodeWidth())
							      .attr("text-anchor", "start");							

						};

					if (!appcompany) { $scope.error = 'Captured data for ' + app + ' is in old data format without company field'; }
					if (!hosts[$scope.app]) { $scope.error = 'No hosts known for app'; }

					$scope.size = (l) => _.keys(l).length;
					// $scope.$watch('app', () => { if (app) { recompute(); } });
					recompute();
					console.info('sankey');

					$scope.hosts = hosts;
					$scope.data = data;
					$scope.pitypes = pitypes;
					$scope.details = details;
					window._s = $scope;
				}
		}); // controller
	});