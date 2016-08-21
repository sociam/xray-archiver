/* global angular, _, jQuery, Backbone */

angular.module('dci')
	.config(function ($stateProvider, $urlRouterProvider) {
		$stateProvider.state('experiment', {
		  	url: '/experiment&id',
		  	template:'<div class="experiment">{{ nTasks }} tasks<div ui-view></div></div>',
		  	resolve: {
				pitypes:($http) => $http.get('../mitm_out/pi_by_host.json').then((x) => x.data),
				hosts: ($http) => $http.get('../mitm_out/host_by_app.json').then((x) => x.data),
				details: ($http) => $http.get('../mitm_out/company_details.json').then((x) => x.data),
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, $stateParams, $state, pitypes, hosts, data, details, storage, utils) {
				$scope.data = data;
				$scope.pitypes = pitypes;
				$scope.hosts = hosts;
				$scope.details = details;				
				var apps = $scope.apps = _.uniq(data.map((x) => x.app)),
					range = $scope.range = utils.range;

				console.log('experiment - loaded ', apps.length, 'apps');

				$scope.thisname = 'experiment';
				$scope.load = function(eId) { 	
					// sets current experiment to be this experiment
					utils.assert(eId, "eID not defined" + eId);
					return storage.get(eId).then((doc) => { 
						$scope.experiment = doc; 
						return $scope.experiment; 
					});
				};
				$scope.save = function() { 	
					// saves $scope.experiment to disk --
					if (!$scope.experiment) { throw new Error("experiment not defined"); }
					if (!$scope.experiment.id) { throw new Error("experiment id not defined"); }
					window.exp = $scope.experiment;
					var id = $scope.experiment.id;
					return storage.db.get(id).then((doc) => {
						console.info("Document found, updating");
						_(doc).extend($scope.experiment);
						return storage.db.put(doc);
					}).catch((e) => {
						// document not found
						console.info("Document not found ");
						return storage.db.put(doc);
					}).finally((x) => {
						console.info('finally done experiment ', id);
						return x;
					});
				};

				// if run then we go and run otherwise manage
				if ($stateParams.id) { 
					$state.go('experiment.run', {id:id});
				} else {
					$state.go('experiment.manage');
				}
			}			
		});

		// manage experimental data: view, delete, export
		$stateProvider.state('experiment.manage', {
		  	url: '/manage',
		  	templateUrl:'tmpl/manage.html',
			controller:function($scope, $state, storage, utils) {
				$scope.selected = 0;
				// load all things into a view
				storage.allDocs().then((docs) => {
					console.info('got all docs ', docs.length);
					$scope.experiments = docs;
				});

				$scope.load = (e) => { 
					$state.go('experiment.run', { id: e._id });
				};
				$scope.delete = (e) => { 
					return storage.db.get(e._id).then((do) => storage.db.remove(doc));
				};
				$scope.export = (e) => { 
					$scope.exportData = JSON.stringify(e);	
				};
			}
		});

		$stateProvider.state('experiment.run', {
		  	url: '/run',
		  	templateUrl:'tmpl/run.html',
			controller:function($scope, $state, $stateParams, utils) {
				if (!$stateParams.id) { 
					$state.go('experiment.manage');
					return;
				}
				$scope.load($stateParams.id).then((experiment) => { $scope.e = experiment; });
				$scope.selectApps = () => $state.go('experiment.runselect');
			}
		});

		$stateProvider.state('experiment.runselect', {
		  	url: '/runselect',
		  	templateUrl:'tmpl/run-select.html',
			controller:function($scope, $state, utils) {
				// allow people to select apps they're using
				$scope.selected = {};
				if ($scope.experiment) { $state.go('experiment.manage'); }
				$scope.save = () => {
					$scope.experiment.pdciApps = _($scope.selected).pickBy((v,k) => v).keys();
					$scope.save().then(() => {
						console.info('experiment saved >> ');
						$state.go('experiment.run', {id:$scope.experiment._id});
					});
				};
			}
		});
		$stateProvider.state('experiment.runtask', {
		  	url: '/runtask',
		  	templateUrl:'tmpl/run-task.html',
			controller:function($scope, $state, utils) {
			}
		});

		$stateProvider.state('experiment.config', {
		  	url: '/config',
		  	templateUrl:'tmpl/config.html',
			controller:function($scope, $state, utils) {
				var hosts = $scope.hosts, data = $scope.data, pitypes = $scope.pitypes;
				console.log('config');
				var apps = $scope.apps,
					range = $scope.range,
					ifaces = $scope.ifaces = ['table','sankey','box'];
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
							$scope.rounds.dci.push({a:apps[0],b:apps[0],c:ifaces[0]});
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

				$scope.save = () => {
					console.info('setting parent');
					$scope.experiment = { 
						_id: $scope.runid,
						participant:$scope.participantid,
						rounds: $scope.rounds,
						configured:true
					};
					$scope.save();
					$state.go('experiment.manage');
				};

				window._s = $scope;
			}			
		});
	});
