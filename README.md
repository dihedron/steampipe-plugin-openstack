# steampipe-plugin-openstack

A Steampipe plugin to query OpenStack data.

```sql
select * from openstack_instance where id = 'foo';
```

Run as:

```bash
$> steampipe query "select * from openstack_instance where id = 'foo';"
```