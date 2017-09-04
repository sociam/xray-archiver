begin;

create table apps(
  id        text primary key not null,
  versions int[]
);

create table app_versions(
  id                      serial         primary key not null,
  app                       text references apps(id) not null,
  store                     text                     not null,
  region                    text                     not null,
  version                   text                     not null,
  screen_flags               int                             ,
  downloaded                bool                     not null,
  analyzed                  bool                     not null,
  last_dl_attempt      timestamp                             ,
  icon                      text                             ,
  uses_reflect              bool                             ,
  last_analyze_attempt timestamp                             ,
  last_alt_checked     timestamp
);

create table developers(
  id         serial primary key not null,
  email      text[]             not null,
  name         text             not null,
  store_site   text                     ,
  site         text
);

create table alt_apps(
   app_id                text             references apps(id) not null,
   alt_app_title         text                                 not null,
   alt_to_url            text                                 not null,
   g_play_url            text                                         ,
   g_play_id             text                                         ,
   icon_url              text                                         ,
   official_site_url     text                                         ,
   is_scraped            bool                                 not null,
   primary key (app_id, alt_app_title)
);

create table playstore_apps(
  id                      int primary key references app_versions(id) not null,
  title                  text                                         not null,
  summary                text                                                 ,
  description            text                                                 ,
  store_url              text                                         not null,
  price                  text                                         not null,
  free                   bool                                         not null,
  rating         numeric(2,1)                                                 ,
  num_reviews          bigint                                                 ,
  genre                  text                                                 ,
  family_genre           text                                                 ,
  min_installs         bigint                                                 ,
  max_installs         bigint                                                 ,
  developer               int               references developers(id) not null,
  updated                date                                         not null,
  android_ver            text                                         not null,
  content_rating         text                                                 ,
  screenshots          text[]                                                 ,
  video                  text                                                 ,
  recent_changes       text[]                                                 ,
  crawl_date             date                                         not null,
  permissions          text[]
);

create table search_terms(
  search_term            text                              primary key not null,
  last_searched          date                                          not null
);

create table app_packages (
  id            int  references app_versions(id) primary key not null,
  packages   text[]                                          not null
);

create table app_perms(
  id             int references app_versions(id) primary key not null,
  permissions text[]                                         not null
);

-- Contains the hostnames that were found in apps via analysis
create table app_hosts(
  id       int references app_versions(id) primary key not null,
  hosts text[]                                                 ,
  pis    int[]
);

create table manual_alts(
  source_id  text not null,
  alt_id     text not null,
  primary key (source_id, alt_id)
);

create table companies(
  id             text     primary key not null,
  name           text                 not null,
  hosts         int[]                         ,
  founded        text                         ,
  acquired       text                         ,
  type         text[]                         ,
  typetag        text                         ,
  jurisdiction   text                         ,
  parent         text references companies(id),
  capital        text                         ,
  equity         text                         ,
  min_size        int                         ,
  max_size        int                         ,
  data_sources text[]                         ,
  description    text
);

create table hosts(
  hostname text     primary key not null,
  company  text references companies(id)
);

create table company_domains (
  company text not null,
  domain text not null,
  primary key(company, domain)
);

create user explorer;
create user retriever;
create user downloader;
create user analyzer;
create user apiserv;
create user suggester;

grant insert, select on search_terms to explorer;

grant select, insert, update on apps to retriever;
grant select, insert on app_versions to retriever;
grant usage on app_versions_id_seq to retriever;
grant select, insert on playstore_apps to retriever;
grant select, update on search_terms to retriever;
grant select, insert, update on developers to retriever;
grant usage on developers_id_seq to retriever;

grant select, update on app_versions to downloader;
grant select on playstore_apps to downloader;

grant select, insert, update on apps to analyzer;
grant select, insert, update on app_versions to analyzer;
grant usage on app_versions_id_seq to analyzer;
grant select  on playstore_apps to analyzer;
grant select, insert, update on app_perms to analyzer;

grant select, insert, update on app_hosts to analyzer;
grant select on companies to analyzer;
grant select, insert, update on alt_apps to analyzer;

grant select on apps to apiserv;
grant select on app_versions to apiserv;
grant select on playstore_apps to apiserv;
grant select on developers to apiserv;
grant select on app_perms to apiserv;
grant select on app_hosts to apiserv;
grant select on companies to apiserv;
grant select on hosts to apiserv;
grant select on alt_apps to apiserv;
grant select on manual_alts to apiserv;
grant select on company_domains to apiserv;

grant select, update, insert on alt_apps to suggester;
grant select, update, insert on manual_alts to suggester;
grant select, update, insert on company_domains to suggester;
grant select, update on app_versions to suggester;
grant select on playstore_apps to suggester;

-- SUMMARY STATISTICS VIEWS
create view all_hosts as
  select unnest(ah.hosts) as hosts from app_hosts ah; 
  
create view host_freq as
  select un.hosts as host_name, bigcnt.big_n, count(un.hosts) as little_n, count(un.hosts)/bigcnt.big_n::float as n_pct from 
    (select count(hosts) as big_n from all_hosts) as bigcnt,
    all_hosts as un
      group by hosts, big_n
      order by little_n using >;

commit;

