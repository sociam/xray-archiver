
begin;


---------------------------------------------------------------------------------------------------
-- MATERIALISED VIEWS TO INCREASE SPEED 
-- GOTTA GO FAST 🦔
---------------------------------------------------------------------------------------------------

--Setup materilized view namespace
create schema matview

--Ensure extension for index matching is setup
create extension pg_trgm from unpackaged;

-----
-- Access speed for playstore apps data
-----
drop table if exists matview.apps_play_data cascade;

create materilized view matview.apps_play_data as select 
a.id, a.title, a.summary, a.description, a.store_url, a.price, a.free, a.rating, 
	a.num_reviews, a.genre, a.family_genre, a.min_installs, a.max_installs, a.updated, 
	a.android_ver, a.content_rating, a.recent_changes, v.app, v.store, v.region,
	v.version, v.icon, v.analyzed, d.email, d.name, d.store_site, d.site, h.hosts, p.permissions,
	pkg.packages from playstore_apps a FULL OUTER join app_versions v on (a.id = v.id) 
	FULL OUTER join developers d on (a.developer = d.id)  
	FULL OUTER join app_hosts h on (a.id = h.id) 
	FULL OUTER join app_perms p on (a.id = p.id) 
	FULL OUTER join app_packages pkg  on (a.id = pkg.id)  order by max_installs;

create index apps_title_gin_trgm_idx on matview.apps_play_data USING GIN (title gin_trgm_ops);

grant select on matview.apps_play_data to apiserv;



commit;