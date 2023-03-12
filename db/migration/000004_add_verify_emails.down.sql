drop table if exists "verify_emails" cascade;

alter table users drop column "is_email_verified";
