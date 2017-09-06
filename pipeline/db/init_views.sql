begin;

-- SUMMARY STATISTICS VIEWS

-- Average number of hosts for each genre
create view genre_host_averages as
  select category, host_count, app_count, host_count/app_count::float as host_avg from (
    select p.genre as category, sum(array_length(v.hosts,1)
  ) as host_count, count(p.genre) as app_count
  from playstore_apps p inner join app_hosts v on (p.id = v.id)
    group by genre) as genre_freq;

grant select on genre_host_averages to apiserv;


create table distinct_hosts as
  select distinct hosts from ( 
    select unnest(hosts) as hosts from app_hosts
  ) as unpack_hosts;

create view distinct_app_hosts as
  select distinct id, hosts from ( 
    select id, unnest(hosts) as hosts from app_hosts
  ) as unpack_hosts;


create view host_app_coverage as


-- Table of all possible Host names and a heuristic regex for the domain of the host.
-- create view host_domains as
--   select hosts,
--          substring(hosts from '(([^\.]*)\.([^\.]*)$)') as domain,
--          substring(hosts from '(([^\.]*)\.([^\.]*)\.([^\.]*)$)') as domain_plus from ( 
--     select distinct hosts from ( 
--       select unnest(hosts) as hosts from app_hosts
--     ) as unpack_hosts
--   ) as distinct_hosts;

grant select on host_domains to apiserv

-- Table of Host, heuristic based domain and company for that domain.
-- create view host_domain_companies as 
--   select distinct d.hosts, d.domain, d.domain_plus, coalesce(c.company, 'unknown') as company
--     from host_domains d, company_domains c
--       where d.domain = c.domain 
--          or d.domain_plus = c.domain;

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
--     (select count(hosts) as big_n from all_hosts)1 as bigcnt,
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

commit;