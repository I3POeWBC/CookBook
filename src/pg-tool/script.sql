create user lyak with password 'xx1234'
/

create table one (
  oneid integer primary key,
  onename text null
)
/

grant select on one to lyak