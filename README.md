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
- [X] Create connection to OpenStack APIs
- [X] Create connection to Compute APIs
- [X] Implement "get VM instance" (by ID)
- [ ] Fill all fields from VM instance
- [ ] Implement "list VM instances"
- [ ] Filter by project ID
- [ ] Filter by project name
- [ ] Filter by hypervisor
- [ ] Filter by availability zone
- [ ] More filter criteria
- [ ] Create connection to Identity APIs
- [ ] Implement "get project" (by ID)
- [ ] Check that joins between instance and project (by ID) work
- [ ] ...
