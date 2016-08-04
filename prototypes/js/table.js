angular.module('dci')
	.service('utils', { 
		return {
			
		};
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
					id2names = $scope.id2names = _.keys(details).reduce((a,id) => { 
						a[id] = details[id].company; return a; 
					}, {}),
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
					checkSize = () => {},
					matchCompany = (x) => appcompany && ((x || '').toLowerCase() === appcompany.toLowerCase()),
					isType = $scope.isType = (id, type) => id && 
						details[id] && 
						details[id].typetag && 
						details[id].typetag.indexOf(type) >= 0,
					getName = $scope.getName = (id) => details[id] && details[id].company || id,
					is3rdPartyType = $scope.is3rdPartyType = (id, type) => isType(id,type) &&
							!_.some([id2names[id], id].map(matchCompany)),
					recompute = () => {
						var apphosts = _(hosts[$scope.app]).pickBy((val) => val > $scope.threshold).keys().value();
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
						$scope.categories = $scope.categories = {
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