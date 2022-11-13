# steampipe-plugin-openstack

A Steampipe plugin to query OpenStack data.

```sql
select * from openstack_instance where id = 'foo';
```

Run as:

```bash
$> steampipe query "select * from openstack_instance where id = 'foo';"
```

# TODO

This plugin is still in the very early stages.

- [x] Skeleton
- [x] Configuration schema
- [ ] Create connection to OpenStack APIs
- [ ] Implement VM instance get (by ID)
- [ ] Implement VM instance list
- [ ] Filter by project ID
- [ ] Filter by project name
- [ ] Filter by hypervisor
- [ ] Filter by availability zone
- [ ] ...
