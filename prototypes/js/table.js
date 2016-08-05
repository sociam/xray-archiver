angular.module('dci')
	.service('utils', { 
		var utils = {
			makeId2Names:(details) => _.keys(details).reduce((a,id) => { a[id] = details[id].company; return a; }, {}),
			makeHTC:(data) => data.reduce((r,a) => {
				if (a.host_company) { 
					r[a.host] = a.host_company; 
					r[a.host_2ld] = a.host_company;
				}
				return r;
			}, {}),
			makeHTH:(data) => data.reduce((r,a) => {
				// host -> 2ld
				if (a.host_2ld) { 
					r[a.host] = a.host_2ld;
				} else { console.error('warning no 2ld ', a.host); }
				return r;
			}, {}),
			isType: (details, id, type) => id && 
				details[id] && 
				details[id].typetag && 
				details[id].typetag.indexOf(type) >= 0,
			matchCompany : (appcompany, x) => appcompany && 
				((x || '').toLowerCase() === appcompany.toLowerCase()),
			is3rdPartyType: (details, id, type) => utils.isType(details,id,type) && 
				!_.some([id2names[id], id].map(utils.matchCompany)),
			makeCompany2pi: (data, hosts, pitypes, threshold) => {
				var apphosts = _(hosts[$scope.app]).pickBy((val) => val > threshold || 0).keys().value(),
					hTh = utils.makeHTH(data),
					hTc = utils.makeHTC(data);
				// next we wanna group together all the pi_types, and consolidate around company
				// console.info('threshold', $scope.threshold, 'apphosts', apphosts.length);
				return apphosts.reduce((r,host) => {
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
			},
			makeCategories:(appCompany, details, c2pi) => {
				'app-publisher': _.pickBy(c2pi, (pis, company) => 
					utils.matchCompany(appCompany, company)),
				'app-functionality': _.pickBy($scope.company2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(details, company'app')),
				'marketing': _.pickBy($scope.company2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(details, company'marketing')),
				'usage tracking': _.pickBy($scope.company2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(details, company 'usage')),
				'payments':_.pickBy($scope.company2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(details, company 'payments')),
				'security':_.pickBy($scope.company2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(details, company 'security')),
				'other': _.pickBy($scope.company2pi, (pis, company) => 
					!utils.matchCompany(appCompany, company) &&							
					!utils.isType(details, company, 'ignore') && 							
					!utils.isType(details, company, 'app') && 							
					!utils.isType(details, company,'marketing') &&
					!utils.isType(details, company, 'platform') && 							
					!utils.isType(details, company, 'usage') &&
					!utils.isType(details, company, 'payments') &&
					!utils.isType(details, company, 'security'))
			}
		};
		return utils;
	}).config(function ($stateProvider, $urlRouterProvider) {
		$stateProvider.state('dci.table', {
			url: '/table',
			templateUrl: 'tmpl/table.html',
			resolve: {
				pitypes:($http) => $http.get('../mitm_out/pi_by_host.json').then((x) => x.data),
				hosts: ($http) => $http.get('../mitm_out/host_by_app.json').then((x) => x.data),
				details: ($http) => $http.get('../mitm_out/company_details.json').then((x) => x.data),
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, pitypes, hosts, details, data, utils, $stateParams) {
				data = $scope.data = data.filter((x) => x.app === $stateParams.app);

				var app = $scope.app = $stateParams.app,
					appcompany = $scope.appcompany = data[0].company,
					id2names = $scope.id2names = utils.makeId2Names(details),
					isType = $scope.isType = (id,type) => utils.isType(details,id,type),
					matchCompany = (x) => utils.matchCompany(appcompany, x)
					getName = $scope.getName = (id) => details[id] && details[id].company || id,
					is3rdPartyType = $scope.is3rdPartyType = (id,type) => utils.is3rdPartyType(details,id,type),
					c2pi = utils.makeCompany2pi(data, hosts),
					cat2c2pi = utils.makeCategories(c2pi);

					recompute = () => {					
						// each of the boxes
						$scope.company2pi = c2pi;
						$scope.categories = $scope.categories = 
					};

				if (!appcompany) { $scope.error = 'Captured data for ' + app + ' is in old data format without company field'; }
				if (!hosts[$scope.app]) { $scope.error = 'No hosts known for app'; }

				$scope.size = (l) => _.keys(l).length;
				$scope.threshold = 0;
				$scope.$watch('threshold', () => { if ($scope.threshold!==undefined) { recompute(); }});

				$scope.hosts = hosts;
				$scope.data = data;
				$scope.pitypes = pitypes;
				$scope.details = details;
				window._s = $scope;
			}
		});