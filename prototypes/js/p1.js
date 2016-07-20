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
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, pitypes, hosts, data, $stateParams) {
				console.log('boxdci stateparams', $stateParams);
				console.log('got relevant ', data.length);
				
				data = data.filter((x) => x.app == $stateParams.app);
				var app = $scope.app = $stateParams.app,
					appCompany = $scope.appcompany = data[0].company,
					hTc = $scope.hTc = data.reduce((r,a) => {
						if (a.host_company) { 
							r[a.host] = a.host_company; r[a.host_2ld] = a.host_company;
						}
						return r;
					}, {}),
					recompute = () => {
						var apphosts = _(hosts[$scope.app]).pickBy((val) => val > $scope.threshold).keys().value();
						// next we wanna group together all the pi_types, and consolidate around company
						console.info('threshold', $scope.threshold, 'apphosts', apphosts.length);
						$scope.company2pi = apphosts.reduce((r,host) => {
							var company = hTc[host],
								host_pis = pitypes[host] || [];
							if (!company) { console.error('no company for host ', host); return r; }
							r[company] = _.union(r[company] || [], host_pis);
							return r;
						}, {});
					};

				if (!hosts[$scope.app]) { $scope.error = 'No hosts known for app'; }

				$scope.threshold = 2;
				$scope.$watch('threshold', () => { if ($scope.threshold!==undefined) { recompute(); }});

				$scope.hosts = hosts;
				$scope.data = data;
				$scope.pitypes = pitypes;
				window._s = $scope;
			}
		});
	});