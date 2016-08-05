/* global angular, _, jQuery, Backbone */
console.log('table js');
angular.module('dci')
	.config(function ($stateProvider, $urlRouterProvider) {
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
				console.info('hello table');
				data = $scope.data = data.filter((x) => x.app === $stateParams.app);
				var app = $scope.app = $stateParams.app,
					appcompany = $scope.appcompany = data[0].company,
					getName = $scope.getName = (id) => details[id] && details[id].company || id,
					c2pi = $scope.c2pi = utils.makeCompany2pi(app, data, hosts, pitypes, 0),
					cat2c2pi = $scope.cat2c2pi = utils.makeCategories(appcompany, details, c2pi),
					recompute = () => {					
						// each of the boxes
						$scope.pilabels = utils.pilabels;
						$scope.pitypes = _(c2pi).values().flatten().uniq().sort((x) => utils.pilabels[x]).value();
						$scope.companies = _(c2pi).keys().uniq().sort((x) => getName(x)).value();
					};

				if (!appcompany) { $scope.error = 'Captured data for ' + app + ' is in old data format without company field'; }
				if (!hosts[$scope.app]) { $scope.error = 'No hosts known for app'; }
				$scope.size = (l) => _.keys(l).length;
				recompute();

				$scope.hosts = hosts;
				$scope.data = data;
				// $scope.pitypes = pitypes;
				// $scope.details = details;
				window._s = $scope;
			}
		});
	});