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

				var pdciData;
				if ($scope.pdciApps) { 
					pdciData = data.filter((x) => $scope.pdciApps.indexOf(x.app) >= 0);
					console.info('pdci data got rows ', pdciData.length);
				}

				data = $scope.data = data.filter((x) => x.app === $stateParams.app);

				$scope.toPairs = (o) => _.toPairs(o).map((x) => { return { key:x[0], val:x[1] }; });

				var app = $scope.app = $stateParams.app,
					appcompany = $scope.appcompany = data[0].company,
					getName = $scope.getName = (id) => details[id] && details[id].company || id,
					c2pi = $scope.c2pi = utils.makeCompany2pi(app, data, hosts, pitypes, 0),
					pdcic2pi = $scope.pdcic2pi = $scope.pdciApps && $scope.pdciApps.length ? utils.makePDCIc2pi($scope.pdciApps, pdciData, hosts, pitypes, 0) : {},
					cat2c2pi = $scope.cat2c2pi = utils.makeCategories(appcompany, details, c2pi),
					recompute = () => {					
						// each of the boxes
						$scope.pilabels = utils.pilabels;
						$scope.pitypes = _(c2pi).values().flatten().uniq().sort((x) => utils.pilabels[x]).value();

						// make nested structure for ng-repeat
						$scope.companies = _(cat2c2pi).keys().map((catname) => {
							return _(cat2c2pi[catname]).keys().map((cn) => {
								return { 
									id: cn,
									category:catname,
									pitypes: cat2c2pi[catname][cn]
								};
							}).value();
						}).flatten().value();
						// $scope.companies = _(c2pi).keys().uniq().sort((x) => getName(x)).value();
					};

				$scope.numCompanies = (cat) => ($scope.cat2c2pi[cat] && _.keys($scope.cat2c2pi[cat]).length) || 0;

				if (!appcompany) { $scope.error = 'Captured data for ' + app + ' is in old data format without company field'; }
				if (!hosts[$scope.app]) { $scope.error = 'No hosts known for app'; }
				$scope.size = (l) => _.keys(l).length;
				recompute();
				$scope.isCat = (c,cat) => cat2c2pi[cat] && cat2c2pi[cat][c] !== undefined;
				$scope.hosts = hosts;
				$scope.data = data;
				// $scope.pitypes = pitypes;
				// $scope.details = details;
				window._ss = $scope;
			}
		});
	});