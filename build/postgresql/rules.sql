create table rules
(
    name   varchar(18) not null,
    subnet cidr        not null
        constraint rules_uniq_subnet unique,
    type   varchar(5)  not null
        constraint rules_type_check check ((type)::text IN ('white', 'black'))
);

alter table rules
    owner to tanker;

