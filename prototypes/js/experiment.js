/* global angular, _, jQuery, Backbone */

angular.module('dci')
	.config(function ($stateProvider, $urlRouterProvider) {
		console.log('... define experiment states');

		$stateProvider.state('experiment', {
		  	url: '/experiment?id',
		  	template:'<div ui-view></div>',
		  	resolve: {
				pitypes:($http) => $http.get('../mitm_out/pi_by_host.json').then((x) => x.data),
				hosts: ($http) => $http.get('../mitm_out/host_by_app.json').then((x) => x.data),
				details: ($http) => $http.get('../mitm_out/company_details.json').then((x) => x.data),
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, $stateParams, $state, pitypes, hosts, data, details, storage, utils) {

				var apps = $scope.apps = _.uniq(data.map((x) => x.app)),
					range = $scope.range = utils.range;				

				console.log('experiment framework > ', apps.length, 'apps');				
				$scope.data = $scope.allData = data;
				$scope.pitypes = pitypes;
				$scope.hosts = hosts;
				$scope.details = details;

				$scope.thisname = 'experiment';
				$scope.load = function(eId) { 	
					// sets current experiment to be this experiment
					utils.assert(eId, "eID not defined" + eId);
					return storage.db.get(eId).then((doc) => { 
						$scope.experiment = doc; 
						return $scope.experiment; 
					});
				};
				$scope.setExperiment = function(e) { $scope.experiment = e; };
				$scope.save = function() { 	
					// saves $scope.experiment to disk --
					console.log('saving starts >>');
					if (!$scope.experiment) { throw new Error("experiment not defined"); }
					if (!$scope.experiment._id) { throw new Error("experiment id not defined"); }
					window.exp = $scope.experiment;
					var id = $scope.experiment._id;
					return storage.db.get(id).then((doc) => {
						doc = _.extend({}, doc, utils.deAngular($scope.experiment));
						console.info("Document found, updating ", doc);
						return storage.db.put(doc).then(() => doc);
					}).catch(() => {
						// document not found
						var doc = utils.deAngular($scope.experiment);
						console.info("Document not found ", doc);
						return storage.db.put(doc);
					}).finally((x) => {
						console.info('finally done experiment ', id);
						return x;
					});
				};
				$scope.getTask = (taskid) => {
					var tasksplit = taskid.split('::'),
						pdciopt = tasksplit[0],
						round = parseInt(tasksplit[1]);

					console.info('getTask choosing ', pdciopt, ' - round#', round);
					return $scope.experiment.rounds[pdciopt][round];
				};
				$scope.makeTaskId = (task) => {
					var dciidx = $scope.experiment.rounds.dci.indexOf(task);
					if (dciidx >= 0) { 
						return ['dci',''+dciidx].join('::');
					}
					dciidx = $scope.experiment.rounds.pdci.indexOf(task);
					return ['pdci',''+dciidx].join('::');
				};

				// delegate to a substate ::
				//   if run then we go and run otherwise manage
				if ($stateParams.id) { 
					$state.go('experiment.run', {id:$stateParams.id});
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
				var refresh_docs = () => {
					storage.db.allDocs({include_docs:true}).then((docs_arr) => {
						console.info('got all docs ', docs_arr.rows, docs_arr.rows.length);
						$scope.experiments = docs_arr.rows.map((x) => x.doc) || [];
					});
				};
				$scope.toDateStr = (ds) => new Date(parseInt(ds)).toDateString() + ' ' + new Date(parseInt(ds)).toLocaleTimeString();
				$scope.load = (e) => { 	$state.go('experiment.run', { id: e._id }); };
				$scope.delete = (e) => { 
					console.info('delete ', e);
					storage.db.get(e._id).then((doc) => storage.db.remove(doc)).then(refresh_docs);
				};
				$scope.export = (e) => { 
					$scope.exportData = JSON.stringify(e, null, 2);	
				};
				$scope.new = () => {
					$scope.setExperiment({});
					$state.go('experiment.config');
				};
				refresh_docs();
			}
		});

		// config a new one
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
				$scope.participantid = 'part-'+utils.guid(4);
				$scope.runid = 'run-'+utils.guid(4);
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
							$scope.rounds.dci.push({a:apps[0],b:apps[0],cond:ifaces[0]});
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
							$scope.rounds.pdci.push({a:apps[0],b:apps[0],cond:ifaces[0],pdci:true});
						});
						console.info(' >> new PDCI rounds length ', $scope.rounds.pdci.length);
					}
				});				
				// 
				if (!$scope.nDciRounds) { $scope.nDciRounds = 3; }
				if (!$scope.nPDciRounds) { $scope.nPDciRounds = 3; }

				$scope.doSave = () => {
					try {
						console.info('setting parent -> ', $scope.experiment);
						_.extend($scope.experiment, { 
							_id: $scope.runid,
							participant:$scope.participantid,
							rounds: $scope.rounds,
							created: new Date().valueOf(),
							configured:true
						});

						console.info('done setting parent -> ', $scope.experiment);

						$scope.save().then(() => {
							console.info('save() :: success saving, now state go');
							$state.go('experiment.manage');
						}).catch((e) => {
							console.error('save() :: exception thrown ', e);
						});
					} catch(e) { 
						console.log('error saving ');
						console.error(e);
					}
				};
				window._s = $scope;
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
				if (!$scope.experiment || $scope.experiment._id !== $stateParams.id) { 
					$scope.load($stateParams.id);
				}
				$scope.selectApps = () => $state.go('experiment.runselect');
				$scope.run = (task) => {
					var tid = $scope.makeTaskId(task);
					console.info('got task id ', tid);
					$state.go('experiment.runtask', {tid:tid});
				};
			}
		});

		$stateProvider.state('experiment.runselect', {
		  	url: '/runselect',
		  	templateUrl:'tmpl/run-select.html',
			controller:function($scope, $state, $stateParams, utils) {
				// allow people to select apps they're using
				if (!$stateParams.id) { $state.go('experiment.manage'); return;	}
				var init = Promise.resolve();				
				if (!$scope.experiment || $scope.experiment._id !== $stateParams.id) { 
					init = $scope.load($stateParams.id).then(() => console.log('loaded'));
				}
				init.then(() => {
					// now we've loaded
					$scope.selected = ($scope.experiment.pdciApps || []).reduce((obj,k) => { obj[k] = true; return obj; }, {});
					$scope.doSave = () => {
						var selected = _($scope.selected).pickBy((v,k) => v).keys().value();
						console.info('setting selected apps to be ', selected);
						$scope.experiment.pdciApps = selected;
						$scope.save().then(() => {
							console.info('<< experiment saved', $scope.experiment);
							$state.go('experiment.run', {id:$scope.experiment._id});
						});
					};
				});
			}
		});

		$stateProvider.state('experiment.runtask', {
		  	url: '/runtask?tid',
		  	templateUrl:'tmpl/run-task.html',
			controller:function($scope, $state, $stateParams, $interval, utils) {
				if (!$stateParams.id) { $state.go('experiment.manage'); return;	}
				var init = Promise.resolve();				
				if (!$scope.experiment || $scope.experiment._id !== $stateParams.id) { 
					init = $scope.load($stateParams.id);
				}
				init.then(() => {
					var start_time = new Date().valueOf(),
						data = $scope.data,
						task = $scope.t = $scope.getTask($stateParams.tid),
						timer_int = $interval(() => { $scope.elapsed = (new Date()).valueOf() - start_time; });

					console.info("GOT task ", task);
					$scope.companies = {
						a:data.filter((x) => x.app === task.a)[0].appcompany,
						b:data.filter((x) => x.app === task.b)[0].appcompany,
					};

					$scope.choiceMade = (choice) => {
						var end_time = (new Date()).valueOf();
						task.result = { 
							chosen:choice,
							start_time: start_time,
							elapsed: end_time-start_time,
							end_time: end_time,
						};
						$scope.save().then(() => { $scope.chosen = choice; });
					};
					// todo
					$scope.gotoNextTask = () => { $state.go('experiment.manage'); };
					$scope.gotoManage = () => { $state.go('experiment.manage'); };
					$scope.$on('$destroy', () => $interval.cancel(timer_int));
				});
			}
		});
		console.log('yolo');
	});
