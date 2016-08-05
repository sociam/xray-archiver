/* global angular, _, jQuery, Backbone */

angular.module('dci').factory('utils', () => {

	var utils = {
		// company id -> company names
		makeId2Names:(details) => 
			_.keys(details).reduce((a,id) => { 
				a[id] = details[id].company; 
				return a; 
			}, {}),

		// returns { host -> companyid }
		makeHTC:(data) => data.reduce((r,a) => {
			if (a.host_company) { 
				r[a.host] = a.host_company; 
				r[a.host_2ld] = a.host_company;
			}
			return r;
		}, {}),

		// returns { host -> host_2ld }
		makeHTH:(data) => data.reduce((r,a) => {
			// host -> 2ld
			if (a.host_2ld) { 
				r[a.host] = a.host_2ld;
			} else { console.error('warning no 2ld ', a.host); }
			return r;
		}, {}),

		// returns true if company_id is of marketing type t
		isType: (details, id, type) => id && 
			details[id] && 
			details[id].typetag && 
			details[id].typetag.indexOf(type) >= 0,

		// aggressive matching of appcompany 			
		cimatch : (appcompany, x) => appcompany && 
			(x || '').toLowerCase() === appcompany.toLowerCase(),

		// isType and not 1st party
		is3rdPartyType: (appcompany, details, id, type) => 
			utils.isType(details,id,type) && 
			!_.some([details[id] && details[id].company || id, id].map((x) => utils.cimatch(appcompany,x))),

		// compile { company_id -> [pitype1,pitype2] .. } filtering by optional threshold
		makeCompany2pi: (app, data, hosts, pitypes, threshold) => {
			var apphosts = _(hosts[app]).pickBy((val) => val > threshold || 0).keys().value(),
				hTh = utils.makeHTH(data),
				hTc = utils.makeHTC(data);
			// next we wanna group together all the pi_types, and consolidate around company
			// console.info('threshold', $scope.threshold, 'apphosts', apphosts.length);
			return apphosts.reduce((r,host) => {
				var company = hTc[host], 
					host_pis = pitypes[host] || [];
				if (!company) { 
					var mfirst = hTh[host].match(/^([^\.]+)\./);
					if (mfirst) { 
						company = mfirst[1]; 
					} else {	
						console.error('no company for host ', host); return r; 
					}
				} 
				r[company] = _.union(r[company] || [], host_pis);
				return r;
			}, {});
		},
		makeCategories:(appCompany, details, c2pi) => { 
			return {
				'app-publisher': _.pickBy(c2pi, (pis, company) => 
					utils.cimatch(appCompany, company)),
				'app-functionality': _.pickBy(c2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(appCompany, details, company, 'app')),
				'marketing': _.pickBy(c2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(appCompany, details, company, 'marketing')),
				'usage tracking': _.pickBy(c2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(appCompany, details, company, 'usage')),
				'payments':_.pickBy(c2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(appCompany, details, company, 'payments')),
				'security':_.pickBy(c2pi, (pis, company) => 
					!utils.isType(details, company, 'ignore') && 
					!utils.isType(details, company, 'platform') && 
					utils.is3rdPartyType(appCompany, details, company, 'security')),
				'other': _.pickBy(c2pi, (pis, company) => 
					!utils.cimatch(appCompany, company) &&							
					!utils.isType(details, company, 'ignore') && 							
					!utils.isType(details, company, 'app') && 							
					!utils.isType(details, company,'marketing') &&
					!utils.isType(details, company, 'platform') && 							
					!utils.isType(details, company, 'usage') &&
					!utils.isType(details, company, 'payments') &&
					!utils.isType(details, company, 'security'))
			};
		},
		pilabels: {
			USER_PERSONAL_DETAILS: 'personal details',
			USER_LOCATION: 'your location',
			USER_LOCATION_COARSE: 'your approximate location',
			DEVICE_ID:'phone id',
			DEVICE_SOFT:'phone characteristics'
		}
	};
	return utils;		
});