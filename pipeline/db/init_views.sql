begin;
---------------------------------------------------------------------------------------------------
-- SUMMARY STATISTICS VIEWS AND TABLES
---------------------------------------------------------------------------------------------------

---------------------------------------------------------------------------------------------------
-- Average number of hosts for each genre
--
-- Features in the "genre_host_averages" API Endpoint.
---------------------------------------------------------------------------------------------------
drop table if exists genre_host_averages cascade;
create table genre_host_averages as
  select category, host_count, app_count, host_count/app_count::float as host_avg from (
    select p.genre as category, sum(array_length(v.hosts,1)
  ) as host_count, count(p.genre) as app_count
  from playstore_apps p inner join app_hosts v on (p.id = v.id)
    group by genre) as genre_freq;

grant select on genre_host_averages to apiserv;

---------------------------------------------------------------------------------------------------
--  Counts by Genre
--
--  Todo -- API Endpoint.
--
---------------------------------------------------------------------------------------------------
drop table if exists genre_counts;
create table genre_counts as
  select genre, count(*) from playstore_apps p inner join app_versions v on (p.id = v.id) 
    where v.analyzed = true
    group by genre;
grant select on genre_counts to apiserv;

---------------------------------------------------------------------------------------------------
-- Standard Deviation of hosts found in each genre
---------------------------------------------------------------------------------------------------
create view genre_host_variance as
  select squared_sum/n-1 
  from (
    select sum(diff_squared) as squared_sum
  from (
    select ((host_avg - 20.391152325979)*(host_avg-20.391152325979)) as diff_squared
  from genre_host_averages) as sums) as sum ,(
    select count(*) as n
  from genre_host_averages) as count;

---------------------------------------------------------------------------------------------------
-- table of distinct hosts 
---------------------------------------------------------------------------------------------------
drop table if exists distinct_hosts;
create table distinct_hosts as
  select distinct hosts from ( 
    select unnest(hosts) as hosts from app_hosts
  ) as unpack_hosts;

---------------------------------------------------------------------------------------------------
-- table of app-host pairs
---------------------------------------------------------------------------------------------------
drop table if exists distinct_app_hosts;
create table distinct_app_hosts as
  select distinct id, hosts from ( 
    select id, unnest(hosts) as hosts from app_hosts
  ) as unpack_hosts;

---------------------------------------------------------------------------------------------------
-- table counting the number of apps that feature specific hosts.
---------------------------------------------------------------------------------------------------
drop table if exists host_app_coverage;
create table host_app_coverage as
  select hosts, count(*) from distinct_app_hosts
    group by hosts;

grant select on host_app_coverage to apiserv;

---------------------------------------------------------------------------------------------------
-- Table of all possible Host names and a heuristic regex for the domain of the host.
---------------------------------------------------------------------------------------------------
drop table if exists host_domains;
create table host_domains as
  select hosts,
    substring(hosts from '(([^\.]*)\.([^\.]*)$)') as domain,
    substring(hosts from '(([^\.]*)\.([^\.]*)\.([^\.]*)$)') as domain_plus from distinct_hosts;

grant select on host_domains to apiserv;

---------------------------------------------------------------------------------------------------
-- Table of Host, heuristic based domain and company for that domain.
---------------------------------------------------------------------------------------------------
drop table if exists host_domain_companies;
create table host_domain_companies as 
  select distinct d.hosts, d.domain, d.domain_plus, coalesce(c.company, 'unknown') as company, coalesce(c.type, 'unknown') as type
    from host_domains d left outer join company_domains c
      on( d.domain = c.domain 
         or d.domain_plus = c.domain
         or lower(d.hosts) ilike '%' || lower(c.domain) || '%');

grant select on host_domain_companies to apiserv;
---------------------------------------------------------------------------------------------------
-- a mapping of hosts-app pairs to host-company pairs. if an app sends to a company, only
-- marked once.
--
-- NOTE - expand to include genres. would be interesting to see if some types of apps send
-- to different types of companies
-- NOTE - Expand to include information on the 
---------------------------------------------------------------------------------------------------
drop table if exists distinct_app_companies;
create table distinct_app_companies as
  select distinct hdc.company, hdc.type, dah.id from host_domain_companies hdc, distinct_app_hosts dah
  where hdc.hosts = dah.hosts;

---------------------------------------------------------------------------------------------------
--
--  genre company coverage.
--
--
---------------------------------------------------------------------------------------------------

drop table if exists company_genre_coverage;
create table company_genre_coverage as 
select company, company_count, genre, genre_total, (company_count/genre_total::float)*100 as coverage_pct from (
  select  count(dac.company) as company_count, dac.company, p.genre, genre_totals.genre_total from distinct_app_companies dac, playstore_apps p, app_versions v, (
    select count(*) as genre_total, genre from  playstore_apps p, app_versions v
      where v.id  = p.id and v.analyzed = true
      group by genre
  ) as genre_totals
    where dac.id = v.id and dac.id = p.id and v.analyzed = true and genre_totals.genre = p.genre
    group by p.genre, dac.company, genre_totals.genre_total
) as totals;

grant select on company_genre_coverage to apiserv;
---------------------------------------------------------------------------------------------------
-- Counts of the amount of apps that feature a host name tied to a company.
--
-- This table features in the 'app_company_freq' API endpoint.
---------------------------------------------------------------------------------------------------
drop table if exists company_app_coverage;
create table company_app_coverage as
  select distinct t2.company, t2.type, app_count, total_apps, app_count/total_apps::float as company_freq from (
    select company, count(*) as app_count, total_apps
        from distinct_app_companies,( select count(*) as total_apps from app_versions where analyzed = true ) as total_app_count 
          group by company, total_apps
  ) as company_app_counts,
  distinct_app_companies t2
  where company_app_counts.company = t2.company
  order by app_count using >;
 grant select on company_app_coverage to apiserv;

---------------------------------------------------------------------------------------------------
-- Counts of the amount of apps that feature a host name tied to a specific purpose
--
-- This table features in the 'app_type_freq' API endpoint.
---------------------------------------------------------------------------------------------------
 drop table if exists app_type_coverage;
 create table app_type_coverage as 
  select type, app_count, total_apps, app_count/total_apps::float as type_freq from (
    select type, count(*) as app_count, total_apps
        from (select distinct type, id from distinct_app_companies) as distinct_app_types,
             (select count(*) as total_apps from app_versions where analyzed = true ) as total_app_count 
          group by type, total_apps
          order by app_count using >
  ) as type_app_counts;
 grant select on app_type_coverage to apiserv;



-----
-- company per app variance
-----
create table company_per_app_varianvce as
select diffSum/count as variance, diffSum, count from (
  select sum(diff) as diffSum 
  from (
       select ((count-6.98820683086)*(count-6.98820683086)) as diff from (
          select count(*)
          from distinct_app_companies group by id ) as counts ) as diffs) as diffsums,
(select count (*) as count from (select distinct id from distinct_app_companies) as ids
) as counts;


commit;
---------------------------------------------------------------------------------------------------
-- Other Views. might be useful at somepoint. but they turned out to be too slow.
---------------------------------------------------------------------------------------------------

-- SUMMARY STATISTICS VIEWS

-- Average number of hosts for each genre
-- create view genre_host_averages as
--   select category, host_count, app_count, host_count/app_count::float as host_avg
--   from (select p.genre as category, sum(array_length(v.hosts,1)) as host_count, count(p.genre) as app_count
--           from playstore_apps p inner join app_hosts v on (p.id = v.id)
--             group by genre) as genre_freq;

-- grant select on genre_host_averages to apiserv;

-- All Hosts
-- create view all_hosts as
--   select unnest(ah.hosts) as hosts, company from app_hosts ah full outer join company_domains on (hosts = domain ); 

-- Known Hosts and Companies
-- create view known_hosts as 
--   select * from (
--     select hosts, coalesce(company, 'Unknown') as company from 
--       all_hosts full outer join company_domains 
--         on (substring(hosts from '(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?(?:stream\.)?(?:test\.)?(?:ipua\.)?(?:retail\.)?(?:i\.)?(?:idx\.)?(?:snappy\.)?(?:snappy.\.)?(?:ofbiz\.)?(?:post\.)?(?:e\.)?(?:api.map\.)?(?:oauth\.)?(?:map\.)?(?:api.maps\.)?(?:m\.)?(?:\.)?(?:badad\.)?(?:bugs\.)?(?:info\.)?(?:info.static\.)?(?:androidads21\.)?(?:maps\.)?(?:serve\.)?(?:sonuj\.)?(?:sonuj,dev\.)?(?:dev\.)?(?:graph.%s\.)?(?:ws\.)?(?:accounts\.)?(?:account\.)?(?:imp\.)?(?:lh6\.)?(?:%s\.)?(?:hemmabast\.)?(?:cdn.unityads\.)?(?:vid\.)?(?:login\.)?(?:sdk\.)?(?:ssdk\.)?(?:graph\.)?(?:ach\.)?(?:unconf.mobad\.)?(?:live\.)?(?:r\.)?(?:tech\.)?(?:rri\.)?(?:ms\.)?(?:unconf\.)?(?:unrcv\.)?(?:cdn\.)?(?:img\.)?(?:ud\.)?(?:ufs\.)?(?:xml\.)?(?:rt\.)?(?:mads\.)?(?:pdn\.)?(?:settings\.)?(?:cdnjs\.)?(?:assets\.)?(?:market\.)?(?:adwatch\.)?(?:code\.)?(?:a\.)?(?:d\.)?(?:ad\.)?(?:www\.)?(?:googleads\.)?(?:googleads.g\.)?(?:schemas\.)?(?:ads\.)?(?:csi\.)?(?:developer\.)?(?:pro\.)?(?:s3\.)?(?:api\.)?(?:docs\.)?(?:ssl\.)?(?:media\.)?(?:play\.)?(?:plus\.)?(?:pagead2\.)?([^:\/\n]+)') = domain ) 
--     ) as counts where company != 'Unknown';

-- All Host Freq Counts
-- create view host_freq as
--   select un.hosts as host_name, bigcnt.big_n, count(un.hosts) as little_n, count(un.hosts)/bigcnt.big_n::float as n_pct from 

--     (select count(hosts) as big_n from all_hosts) as bigcnt,
--     all_hosts as un
--       group by hosts, big_n
--       order by little_n using >;

-- Known Hosts company Freq Counts
-- create view known_host_freq as
--   select un.company as host_name, bigcnt.big_n, count(un.company) as little_n, count(un.company)/bigcnt.big_n::float as n_pct from 
--     (select count(company) as big_n from known_hosts) as bigcnt,
--     known_hosts as un
--       group by company, big_n
--       order y little_n using >;


-- grant select on all_hosts, known_hosts, host_freq, known_host_freq to apiserv;

