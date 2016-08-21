/* global angular, _, jQuery, Backbone */

angular.module('dci')
	.config(function ($stateProvider, $urlRouterProvider) {
		$stateProvider.state('experiment', {
		  	url: '/experiment',
		  	template:'<div class="experiment">{{ nTasks }} tasks<div ui-view></div></div>',
		  	resolve: {
				pitypes:($http) => $http.get('../mitm_out/pi_by_host.json').then((x) => x.data),
				hosts: ($http) => $http.get('../mitm_out/host_by_app.json').then((x) => x.data),
				details: ($http) => $http.get('../mitm_out/company_details.json').then((x) => x.data),
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, $state, pitypes, hosts, data, utils) {
				console.log('experiment');
				$scope.thisname = 'experiment';
				if (!$scope.configuration) { 
					$state.go('experiment.config');
					return;
				}
				// $scope.nTasks = $scope.configuration.rounds.dci.length + 
				// 	$scope.configuration.rounds.pdci.length;

			}			
		});
		$stateProvider.state('experiment.config', {
		  	url: '/config',
		  	templateUrl:'tmpl/config.html',
		  	resolve: {
				pitypes:($http) => $http.get('../mitm_out/pi_by_host.json').then((x) => x.data),
				hosts: ($http) => $http.get('../mitm_out/host_by_app.json').then((x) => x.data),
				details: ($http) => $http.get('../mitm_out/company_details.json').then((x) => x.data),
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, $state, pitypes, hosts, data, utils) {
				console.log('config');
				var apps = $scope.apps = _.uniq(data.map((x) => x.app)),
					range = $scope.range = utils.range;
				$scope.genID = utils.guid;
				$scope.rounds = { dci: [], pdci: [] };
				$scope.$watch('nDciRounds', () => {
					if (!$scope.nDciRounds || $scope.nDciRounds.length === 0) { return; }

					var pN = parseInt($scope.nDciRounds);
					if (_.isNaN(pN)) { $scope.error('not a valid number ', $scope.nDciRounds); return; }

					if ($scope.rounds.dci.length > pN) { 
						$scope.rounds.dci = $scope.rounds.dci.slice(0,pN);
					} else if ($scope.rounds.dci.length < pN) { 
						console.info(' lengthening ', $scope.rounds.dci.length, ' pN ', pN);
						range(pN - $scope.rounds.dci.length).map(() => {
							$scope.rounds.dci.push({a:apps[0],b:apps[0]});
						});
						console.info(' >> new DCI rounds length ', $scope.rounds.dci.length);
					}
				});
				$scope.$watch('nPDciRounds', () => {
					if (!$scope.nPDciRounds) { return; }
					var pN = parseInt($scope.nPDciRounds);
					if (_.isNaN(pN)) { $scope.error('not a valid number ', $scope.nPDciRounds); return; }
					if ($scope.rounds.pdci.length > pN) { 
						$scope.rounds.pdci = $scope.rounds.pdci.slice(0,pN);
					} else if ($scope.rounds.pdci.length < pN) { 
						console.info(' lengthening ', $scope.rounds.pdci.length, ' pN ', pN);
						range(pN - $scope.rounds.pdci.length).map(() => {
							$scope.rounds.pdci.push({a:apps[0],b:apps[0]});
						});
						console.info(' >> new PDCI rounds length ', $scope.rounds.pdci.length);
					}
				});				
				// 
				if (!$scope.nDciRounds) { $scope.nDciRounds = 3; }
				if (!$scope.nPDciRounds) { $scope.nPDciRounds = 3; }

				$scope.run = () => {
					console.info('setting parent');
					$scope.$parent.configuration = { rounds: $scope.rounds };
					$state.go('experiment');
				};

				window._s = $scope;
			}			
		});
	});
