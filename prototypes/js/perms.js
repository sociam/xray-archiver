/* global angular, _, $ */
angular.module('dci')
	.directive('permissions', function() {
		return {
			templateUrl:'tmpl/perms.html',
			restrict:'E',
			scope:{app:'=', appcompany:'=', 'showCategories':'='},
			controller:function($scope, $timeout, utils) {
				window._ps = $scope;
				var	recompute = () => {					
					var hosts = $scope.$parent.hosts,
						details = $scope.$parent.details,
						pitypes = $scope.$parent.pitypes,
						allData = $scope.$parent.allData,data = $scope.data = allData.filter((x) => x.app === $scope.app),
						c2pi = $scope.c2pi = utils.makeCompany2pi($scope.app, data, hosts, pitypes, 0),
						cat2c2pi = utils.makeCategories($scope.appcompany, details, c2pi),
						pits = $scope.pits = _(c2pi).values().flatten().uniq(),
						pi2cat = $scope.pi2cat = pits.reduce((red, pit) => {
							var cats = _.keys(cat2c2pi).filter((cat) => {
								return _(cat2c2pi[cat]).values().flatten().uniq().indexOf(pit) >= 0;
							});
							red[pit] = cats;
							return red;
						}, {});				
					};

				$scope.$watch('app', () => { if ($scope.app) { recompute(); } });
			}
		};
	});