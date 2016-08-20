/* global angular, _, $ */
angular.module('dci')
	.directive('dciTable', {
		templateUrl: 'tmpl/table.html',
		controller:function($scope, $timeout, utils) {
			var hosts = $scope.hosts,
				details = $scope.details,
				pitypes = $scope.pitypes,
				allData = $scope.allData,
				app = $scope.app,
				appcompany = $scope.appcompany,
				recompute = () => {					
					var data = $scope.data = allData.filter((x) => x.app === $scope.app);
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

				var hide_Timer;
				$scope.showInfoBox = function(data, type, $event) { 
					// figure out x and y 
					// console.info('showinfobox ', $event.target, $($event.target).position());
					if (hide_Timer) { $timeout.cancel(hide_Timer); hide_Timer = undefined; }
					var position = $event && $event.target && $($event.target).position();
					if (type === 'pitype') { 
						$scope.infoboxx = position.left;
						$scope.infoboxy = position.top+60;
						$scope.infobox = { name:data, label: data, type:type };
					} else if (type === 'category') {
						$scope.infoboxx = position.left + 170;
						$scope.infoboxy = position.top;												
						$scope.infobox = { name:data, label: data, type:type };
					} else { 
						console.log('width ', $($event.target).width());
						$scope.infoboxx = position.left + $($event.target).width() + 60;
						$scope.infoboxy = position.top;						
						$scope.infobox = data.details;
					}
					console.log("scope infobox is ", $scope.infobox);
				};
				$scope.hideInfoBox = function() { 
					// console.info('ng mouse leave ', data);
					hide_Timer = $timeout(() => delete $scope.infobox, 100);
				};

				$scope.numCompanies = (cat) => ($scope.cat2c2pi[cat] && _.keys($scope.cat2c2pi[cat]).length) || 0;

				if (!appcompany) { $scope.error = 'Captured data for ' + app + ' is in old data format without company field'; }
				if (!hosts[$scope.app]) { $scope.error = 'No hosts known for app'; }
				$scope.size = (l) => _.keys(l).length;
				recompute();

				$scope.$watch('pdciApps', recompute);
				// $scope.pitypes = pitypes;
				// $scope.details = details;
				window._ss = $scope;
			}
	});
