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
				$scope.apps = _.uniq(data.map((x) => x.app));
			}
		});
		// base dci state
		$stateProvider.state('dci', {
		  	url: '/dci?app',
		  	templateUrl:'tmpl/view.html',
		  	resolve: {
				data: ($http) => $http.get('../mitm_out/data_all.json').then((x) => x.data)
			},
			controller:function($scope, $state, data, $stateParams) {
				$scope.apps = _.uniq(data.map((x) => x.app));
				data = $scope.data = data.filter((x) => x.app === $stateParams.app);				

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
			}
		  });
	}).component('companyInfo', {
	  templateUrl: 'tmpl/company-info.html',
	  replace:true,
	  bindings: { company: '=', companyName:'=', details:'=', showDetails:'=' },	  
	  controller: function($scope) { console.log('company id:', this.company, 'name: ', this.companyName, this.details); }
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
	  bindings: { apps: '=', selected:'=', mode:'=', showCompanyDetails:'=' },
	  controller:function($scope, $state) {
	  	// console.log('selected ', this.selected);
	  	$scope.$watch(() => this.selected + this.mode, () => { 
	  		console.info('new selected app ', this.selected, 'mode: ', this.mode);
	  		if (this.selected && this.mode) { 
	  			var modemap = { box: 'dci.box', sankey: 'dci.sankey', table: 'dci.table' };
		  		console.info('go ', this.selected, this.mode);
		  		$state.go(modemap[this.mode], {app:this.selected}); 
		  	}
	  	});
	  	if (this.showCompanyDetails === undefined) { this.showCompanyDetails = 'hide'; 	}
	  }
   })