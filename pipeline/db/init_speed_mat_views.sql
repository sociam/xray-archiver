
begin;


---------------------------------------------------------------------------------------------------
-- MATERIALISED VIEWS TO INCREASE SPEED 
-- GOTTA GO FAST 🦔
---------------------------------------------------------------------------------------------------

create schema matview


-----
-- Access speed for playstore apps data
-----
drop table if exists matview.apps_play_data cascade;

create MATERIALIZED view matview.apps_play_data as SELECT 
a.id, a.title, a.summary, a.description, a.store_url, a.price, a.free, a.rating, 
	a.num_reviews, a.genre, a.family_genre, a.min_installs, a.max_installs, a.updated, 
	a.android_ver, a.content_rating, a.recent_changes, v.app, v.store, v.region,
	v.version, v.icon, v.analyzed, d.email, d.name, d.store_site, d.site, h.hosts, p.permissions,
	pkg.packages FROM playstore_apps a FULL OUTER JOIN app_versions v ON (a.id = v.id) 
	FULL OUTER JOIN developers d ON (a.developer = d.id)  
	FULL OUTER JOIN app_hosts h ON (a.id = h.id) 
	FULL OUTER JOIN app_perms p ON (a.id = p.id) 
	FULL OUTER JOIN app_packages pkg  ON (a.id = pkg.id)  ORDER BY max_installs;

grant select on matview.apps_play_data to apiserv;



commit;