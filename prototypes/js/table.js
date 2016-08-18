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
				var allData = data;
				$scope.toPairs = (o) => _.toPairs(o).map((x) => { return { key:x[0], val:x[1] }; });
				$scope.pilabels = utils.pilabels;

				var app = $scope.app = $stateParams.app,
					appcompany = $scope.appcompany = data[0].company,
					getName = $scope.getName = (id) => details[id] && details[id].company || id,
					recompute = () => {					

						data = $scope.data = allData.filter((x) => x.app === $stateParams.app);

						var c2pi = $scope.c2pi = utils.makeCompany2pi(app, data, hosts, pitypes, 0),
							cat2c2pi = $scope.cat2c2pi = utils.makeCategories(appcompany, details, c2pi),					
							isPDCI = $scope.pdciApps && $scope.pdciApps.length,
							pdciData = isPDCI && allData.filter((x) => $scope.pdciApps.indexOf(x.app) >= 0),							
							pdcic2pi = $scope.pdcic2pi = isPDCI ? utils.makePDCIc2pi($scope.pdciApps, pdciData, hosts, pitypes, 0) : {},
							pdcicat2c2pi = isPDCI ? utils.makeCategories(appcompany, details, pdcic2pi) : {};

						// each of the boxes
						$scope.pitypes = _(c2pi).values().flatten().union(_(pdcic2pi).values().flatten().value()).uniq().sort((x) => utils.pilabels[x]).value();

						// make nested structure for ng-repeat
						$scope.companies = _(cat2c2pi).keys().map((catname) => {
							return _(cat2c2pi[catname]).keys().filter((cn) => !isPDCI || !pdcic2pi[cn] || !pdcic2pi[cn].length).map((cn) => {
								return { 
									id: cn,
									details: details[cn],
									category:catname,
									pitypes: c2pi[cn]
								};
							}).value();
						}).flatten().value();

						$scope.pdci_companies = isPDCI && _(pdcicat2c2pi).keys().map((catname) => {
							return _(pdcicat2c2pi[catname]).keys().map((cn) => { // .filter((cn) => pdcic2pi[cn] && pdcic2pi[cn].length).
								return { 
									id: cn,
									details: details[cn],									
									relevant_to_app:c2pi[cn] !== undefined,
									category:catname,
									pitypes:pdcic2pi[cn]
								};
							}).value();
						}).flatten().value();
					};

				$scope.showInfoBox = function(data, type, $event) { 
					// figure out x and y 
					// console.info('showinfobox ', $event.target, $($event.target).position());
					var position = $event && $event.target && $($event.target).position();
					$scope.infoboxx = position.left + 20;
					$scope.infoboxy = position.top + 40;
					$scope.infobox = type === 'pitype' ? { label: data, type:'pitype' } : data.details;
					console.log("scope infobox is ", $scope.infobox);
				};

				$scope.numCompanies = (cat) => ($scope.cat2c2pi[cat] && _.keys($scope.cat2c2pi[cat]).length) || 0;

				if (!appcompany) { $scope.error = 'Captured data for ' + app + ' is in old data format without company field'; }
				if (!hosts[$scope.app]) { $scope.error = 'No hosts known for app'; }
				$scope.size = (l) => _.keys(l).length;
				recompute();
				$scope.hosts = hosts;
				$scope.data = data;

				$scope.$watch('pdciApps', recompute);
				// $scope.pitypes = pitypes;
				// $scope.details = details;
				window._ss = $scope;
			}
		});
	});