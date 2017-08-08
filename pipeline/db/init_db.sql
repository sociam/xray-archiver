create table apps(
  id        text primary key not null,
  versions int[]
);

create table app_versions(
  id                 serial         primary key not null,
  app                  text references apps(id) not null,
  store                text                     not null,
  region               text                     not null,
  version              text                     not null,
  screen_flags          int                             ,
  downloaded           bool                     not null,
  analyzed             bool                     not null,
  last_dl_attempt timestamp                     not null,
  icon                 text
);

create table developers(
  id         serial primary key not null,
  email      text[]             not null,
  name         text             not null,
  store_site   text                     ,
  site         text
);

create table alt_apps(
   id            text             references apps(id) not null,
   title         text                                 not null,
   url           text                                         ,
   primary key (id, title)
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

create user explorer;
create user retriever;
create user downloader;
create user analyzer;
create user apiserv;
create user suggester;

grant insert, select on search_terms to explorer;

grant select, insert on apps to retriever;
grant select, insert on app_versions to retriever;
grant select, insert on playstore_apps to retriever;
grant select, update on search_terms to retriever;
grant select, insert, update on developers to retriever;

grant select, update on app_versions to downloader;

grant select, update, insert on apps to analyzer;
grant select, update, insert on app_versions to analyzer;
grant select  on playstore_apps to analyzer;
grant select, insert, update on app_perms to analyzer;
grant select, insert on app_hosts to analyzer;
grant select  on companies to analyzer;

grant select on apps to apiserv;
grant select on app_versions to apiserv;
grant select on playstore_apps to apiserv;
grant select on developers to apiserv;
grant select on app_perms to apiserv;
grant select on app_hosts to apiserv;
grant select on companies to apiserv;
grant select on hosts to apiserv;
grant select on alt_apps to apiserv;

grant select, update, insert on alt_apps to suggester;
grant select on app_versions to suggester;
