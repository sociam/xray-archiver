/* global angular, _, jQuery, Backbone */

angular.module('dci', ['ui.router', 'ngAnimate', 'ngTouch', 'ngSanitize'])
	.controller('p1', function () {})
	.config(function ($stateProvider, $urlRouterProvider) {
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
				$scope.showPDCI = {showing:false};
				$scope.apps = _.uniq(data.map((x) => x.app));
				// $scope.size = (obj) => _.toPairs(obj).length;
			}
		});
		// base dci state
		$stateProvider.state('dci', {
		  	url: '/dci?app&pdciapps',
		  	templateUrl:'tmpl/view.html',
		  	resolve: {
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, $state, data, $stateParams) {
				$scope.apps = _.uniq(data.map((x) => x.app));
				data = $scope.data = data.filter((x) => x.app === $stateParams.app);				

				var refreshpdciApps = () => {
					console.log('updating pdciApps ', $scope.pdciAppsObj);
					$scope.pdciApps = _.keys($scope.pdciAppsObj).filter((k) => $scope.pdciAppsObj[k]);
				};

				if (!data.length) { $scope.error = 'no data for app ' + $stateParams.app; }
				
				// set app name / company
				$scope.app = $stateParams.app;
				$scope.appcompany = data[0].company;

				// sets mode for toolbar
				$scope.mode = {
					'dci.box': 'box',
					'dci.sankey': 'sankey',
					'dci.table' : 'table'
				}[$state.$current.toString()];

				$scope.$watchCollection('pdciAppsObj', refreshpdciApps);

				if ($stateParams.pdciapps && $stateParams.pdciapps.length > 0) {
					var param = typeof $stateParams.pdciapps == 'string' ? [$stateParams.pdciapps] : $stateParams.pdciapps;
					$scope.pdciAppsObj = param.reduce((a, app) => { a[app] = true; return a; }, {});
				}

				window._sD = $scope;
			}
		  });
	}).component('piTypesDisplay', {
	  templateUrl: 'tmpl/pi-types-display.html',
	  bindings: { types: '=' },	  	  
	  controller: function($scope) { 
	  	console.log('pitypes', this.types); 
	  	this.pis = this.types;
	  	var types = this.types;
	  	$scope.contains = ((type) => types.filter((x) => x === type).length);
	  	$scope.containsPartial = ((type) => types.filter((x) => x.indexOf(type) >= 0).length);
	  }
   }).component('toolbar', { 
	  templateUrl: 'tmpl/toolbar.html',
	  bindings: { apps: '=', selected:'=', mode:'=', showCompanyDetails:'=', showpdci:'=', pdciapps:'=' },
	  controller:function($scope, $state) {
	  	// console.log('selected ', this.selected);
	  	$scope.$watch(() => this.selected + this.mode, () => { 
	  		console.info('new selected app ', this.selected, 'mode: ', this.mode, this.pdciapps);
	  		if (this.selected && this.mode) { 
	  			var modemap = { box: 'dci.box', sankey: 'dci.sankey', table: 'dci.table' };
		  		console.info('go ', this.selected, this.mode);
		  		$state.go(modemap[this.mode], {app:this.selected, pdciapps:this.pdciapps}); 
		  	}
	  	});
	  	if (this.showCompanyDetails === undefined) { this.showCompanyDetails = 'hide'; 	}
	  	$scope.$watch(() => this.showPDCI, () => console.log('showPDCI ', this.showPDCI));
	  }
   }).component('pdciAppSelector', {
	  templateUrl: 'tmpl/pdci-app-selector.html',
	  bindings: { apps: '=', selected:'=', showing:"=", pdciapps:'=' },	  	  
	  controller: function($scope) { 
	  	console.log('pdciAppSelector', this.apps, this.selected); 
	  	if (this.selected === undefined) { this.selected = {}; }
	  	// debug crap
	  	// $scope.$watchCollection(() => this.selected, () => { console.log('select watch ! ', this.selected); });
	  	// $scope.$watch(() => this.showing, () => { console.log('showing watch ! ', this.showing); });
	  }
   }).component('companyInfo', { // used by the box display
	  templateUrl: 'tmpl/company-info.html',
	  replace:true,
	  bindings: { company: '=', companyName:'=', details:'=', showDetails:'=' },	  
	  controller: function($scope) { console.log('company id:', this.company, 'name: ', this.companyName, this.details); }
   }).component('companyInfoBox', { // used by sankey and table
	  templateUrl: 'tmpl/company-info-box.html',
	  bindings: { selected: '=', x:'=', y:'=' },	  	  
	  controller: function($scope) { 
		 $scope.close = () => { delete $scope.selected; };
		 var emoji_table = { 
			 US : '&#x1F1FA;&#x1F1F8;',
			 UK : '&#x1F1EC;&#x1F1E7;',
			 AT : '&#x1F1E6;&#x1F1F9;',
			 CN : '&#x1F1E8;&#x1F1F3;',
			 FR : '&#x1F1EB;&#x1F1F7;',
			 CA : '&#x1F1E8;&#x1F1E6;',
			 DE : '&#x1F1E9;&#x1F1EA;'
		};
		 $scope.$watch(() => this.selected, () => { 
		 	var s = $scope.selected = this.selected && _.clone(this.selected) || this.selected;
			if (s && s.company && s.equity && s.equity.length) {
				var n = parseInt(s.equity);
				if (n > 1e6) { s.equity = Math.round(n/1.0e5)/10.0 + "m"; }
				if (n > 1e9) { s.equity = Math.round(n/1.0e8)/10.0 + "bn"; } 
			}		 
			if (s && s.company && s.jurisdiction_code && emoji_table[s.jurisdiction_code.toUpperCase()]) { 
				s.jurisdiction_flag = emoji_table[s.jurisdiction_code.toUpperCase()];
			}
		 });
	  }
   });