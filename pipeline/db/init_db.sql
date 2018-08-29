-----
--
--  Table Creation
--
-----

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
  apk_location              text                             , -- Path to the APK for this version of the App.
  apk_filesystem            text                             ,
  apk_filesystem_name       text                             ,
  apk_location_root         text                             ,
  apk_location_uuid         text                             , -- UUID of the device that this APK is stored on.
  apk_server_location       text                             , -- Really an indicator of what VM the APK is stored on.
  screen_flags               int                             ,
  downloaded                bool                     not null,
  apk_archived              bool                     default false,
  has_apk_stored            bool                     default false,
  analyzed                  bool                     not null,
  last_dl_attempt      timestamp                             ,
  icon                      text                             ,
  uses_reflect              bool                             ,
  last_analyze_attempt timestamp                             ,
  last_alt_checked     timestamp
);

create table ad_hoc_analysis(
  id                    serial        not null primary key,
  app_id                int           not null references app_versions(id),
  analyser_name         text          not null,
  analysis_by           text          not null default 'anon',
  analysis_date         timestamp     not null default now(),
  results               json          not null
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

create table app_companies(
  id         int references app_versions(id) primary key not null,
  companies  text[]
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
  domain  text not null,
  type    text         ,
  primary key(company, domain)
);

--
--    Company Association Information
--

create table iotDevices(
  id                      serial    not null primary key
);

create table websites(
  id                      serial    not null  primary key
);

--
--    Names of companies associated with an App, IoT Device, or website.
--

create table companyNames(
  id                      serial      not null    primary key,
  company_name            text        not null    unique
);

--
--    All apps, IoT devices, and websites associated with a given company.
--

create table companyAssociations(
  id                      serial      not null    primary key,
  company_name            text        not null    unique references companyNames(company_name),
  app_associations        int[],
  iot_device_associations int[],
  website_associations    int[]
);

create table companyAppAssociations(
  id                      serial      not null    ,
  company_name            text        not null    references companyNames(company_name),
  associated_app          serial      not null    references app_versions(id),
  primary key (company_name, associated_app)
);

create table companyIoTDeviceAssociations(
  id                      serial      not null    ,
  company_name            text        not null    references companyNames(company_name),
  associated_iot_device   serial      not null    references iotDevices(id),
  primary key (company_name, associated_iot_device)
);

create table companyWebsiteAssociations(
  id                      serial      not null    ,
  company_name            text        not null    references companyNames(company_name),
  associated_website      serial      not null    references websites(id),
  primary key (company_name, associated_website)
);

commit;

-----
--
--  Function Creation
--
-----
begin;
create or replace function createCompanyAssociationRecord() returns trigger as
  $BODY$
    begin
      insert into companyAssociations(
        company_name,
        app_associations,
        iot_device_associations,
        website_associations
      ) values (
        new.company_name,
        array[]::integer[],
        array[]::integer[],
        array[]::integer[]
      );
      return new;
    end;
  $BODY$
language plpgsql;


create or replace function updateCompanyAppAssociations() returns trigger as
  $BODY$
    begin
      update companyAssociations
        set app_associations = app_associations || new.associated_app
          where company_name = new.company_name;
      return new;
    end;
  $BODY$
language plpgsql;

create or replace function updateCompanyWebsiteAssociations() returns trigger as
  $BODY$
    begin
      update companyAssociations
        set website_associations = website_associations || new.associated_website
          where company_name = new.company_name;
      return new;
    end;
  $BODY$
language plpgsql;

create or replace function updateCompanyIoTDeviceAssociations() returns trigger as
  $BODY$
    begin
      update companyAssociations
        set iot_device_associations = iot_device_associations || new.associated_iot_device
          where company_name = new.company_name;
      return new;
    end;
  $BODY$
language plpgsql;
commit;

-----
--
--  Trigger Creation
--
-----
begin;
create trigger onCompanyNameInsert
  after insert on companyNames
    for each row
      execute procedure createCompanyAssociationRecord();

create trigger onCompanyAppAssociationInsert
  after insert on companyAppAssociations
    for each row
      execute procedure updateCompanyAppAssociations();

create trigger onCompanyWebsiteAssociationInsert
  after insert on companyWebsiteAssociations
    for each row
      execute procedure updateCompanyWebsiteAssociations();

create trigger onCompanyIoTDeviceAssociationInsert
  after insert on companyIoTDeviceAssociations
    for each row
      execute procedure updateCompanyIoTDeviceAssociations();

commit;

-----
--
--  Role Creation
--
-----
begin;

create user explorer;
create user retriever;
create user downloader;
create user analyzer;
create user apiserv;
create user suggester;

commit;

-----
--
--  Role Permissions
--
-----
begin;

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
grant select on apps to downloader;

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

commit;

-- Query to migrate existing analysis into the ad_hoc_analysis, complete with json!
-- insert into ad_hoc_analysis (app_id, analyser_name, analysis_by, results)
--   select  coalesce(hosts_id, perms_id, packages_id),
--           'Golang analyser',
--           'A.D.S Team',
--           row_to_json(q)
--   from (
--     select *,
--       app_hosts.id as hosts_id,
--       app_packages.id as packages_id,
--       app_perms.id as perms_id
--     from
--       app_hosts
--     full outer join
--       app_perms
--     on
--       app_hosts.id = app_perms.id
--     full outer join
--       app_packages
--     on
--       app_hosts.id = app_packages.id
--     and
--       app_perms.id = app_packages.id
--   ) as q;







insert into ad_hoc_analysis (app_id, analyser_name, analysis_by, results)
  select  coalesce(hosts_id, perms_id, packages_id),
          'Golang analyser',
          'A.D.S Team',
          row_to_json(q)
  from (
    select *,
      app_hosts.id as hosts_id,
      app_packages.id as packages_id,
      app_perms.id as perms_id
    from
      app_hosts
    full outer join
      app_perms
    on
      app_hosts.id = app_perms.id
    full outer join
      app_packages
    on
      app_hosts.id = app_packages.id
    and
      app_perms.id = app_packages.id
  ) as q;