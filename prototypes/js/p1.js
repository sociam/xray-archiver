/* global angular, _, jQuery, Backbone */

angular.module('dci', ['ui.router', 'ngAnimate', 'ngTouch', 'ngSanitize'])
	.controller('p1', function () {

	}).config(function ($stateProvider, $urlRouterProvider) {
	    $urlRouterProvider.otherwise('/chooseapp');
		$stateProvider.state('chooseApp', {
			url: '/chooseapp',
			templateUrl: 'tmpl/choose-app.html',
			resolve: {
				pitypes:($http) => $http.get('../mitm_out/pi_by_host.json').then((x) => x.data),
				hosts: ($http) => $http.get('../mitm_out/host_by_app.json').then((x) => x.data),
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, pitypes, hosts, data) {
				console.log('chooseapp');
				window._pit = pitypes;
				window._hosts = hosts;
				window._data = data;
				$scope.pitypes = pitypes;
				$scope.hosts = hosts;
				$scope.data = data;
				$scope.apps = _.uniq(data.map((x) => x.app));
			}
		});
		$stateProvider.state('boxdci', {
			url: '/boxdci?app',
			templateUrl: 'tmpl/box-dci.html',
			resolve: {
				pitypes:($http) => $http.get('../mitm_out/pi_by_host.json').then((x) => x.data),
				hosts: ($http) => $http.get('../mitm_out/host_by_app.json').then((x) => x.data),
				details: ($http) => $http.get('../mitm_out/company_details.json').then((x) => x.data),
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, pitypes, hosts, details, data, $stateParams) {
				console.log('boxdci stateparams', $stateParams);
				console.log('got relevant ', data.length);
				
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
						}
						return r;
					}, {}),
					matchCompany = (x) => (x || '').toLowerCase() === appcompany.toLowerCase(),
					// isAd = $scope.isAd = (id) => {
					// 	console.info(id, ' ', id2names[id], id, 'details ', details[id]);
					// 	return id && 
					// 		!_.some([id2names[id], id].map(matchCompany)) && 
					// 		details[id] && details[id].typetag && details[id].typetag.indexOf('advert') >= 0;
					// },
					is3rdPartyType = $scope.is3rdPartyType = (id, type) => {
						return id && 
							!_.some([id2names[id], id].map(matchCompany)) &&  // filter out self
							details[id] && details[id].typetag && details[id].typetag.indexOf(type) >= 0;
					},
					isAd = (id) => is3rdPartyType(id,'advert'),
					recompute = () => {
						var apphosts = _(hosts[$scope.app]).pickBy((val) => val > $scope.threshold).keys().value();
						// next we wanna group together all the pi_types, and consolidate around company
						console.info('threshold', $scope.threshold, 'apphosts', apphosts.length);
						$scope.company2pi = apphosts.reduce((r,host) => {
							var company = id2names[hTc[host]] || hTc[host], host_pis = pitypes[host] || [];
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
						$scope.appcompany2pi = _.pickBy($scope.company2pi, (pis, company) => matchCompany(company));
						$scope.ad2pi = _.pickBy($scope.company2pi, (pis, company) => isAd(company));
						$scope.analytics2pi = _.pickBy($scope.company2pi, (pis, company) => is3rdPartyType(company, 'analytics'));
						$scope.non2pi = _.pickBy($scope.company2pi, (pis, company) => !$scope.appcompany2pi[company] &&
							!isAd(company) && 
							!is3rdPartyType(company, 'analytics'));
					};

				// $scope.details = companydetails;

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
	});