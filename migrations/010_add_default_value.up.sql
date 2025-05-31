alter table session alter column user_agent set default '';
alter table session alter column ip_address set default '';

alter table conversation alter column user_id set default '';
alter table conversation alter column last_message_id set default '';