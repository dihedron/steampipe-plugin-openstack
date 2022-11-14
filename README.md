# steampipe-plugin-openstack

A Steampipe plugin to query OpenStack data.

```sql
select * from openstack_instance where id = 'foo';
```

Run as:

```bash
$> steampipe query "select * from openstack_instance where id = 'foo';"

$> steampipe query "select vm.id, vm.name, vm.host_id, vm.flavor_sockets, vm.flavor_disk, prj.name, prj.enabled, prj.id from openstack_instance vm, openstack_project prj where vm.id = '12345678-90ab-cdef-1234-567890abcdef' and vm.project_id = prj.id;"
```

# TODO

This plugin is still in the very early stages.

- [x] Skeleton
- [x] Configuration schema
- [X] Create connection to OpenStack APIs
- [X] Create connection to Compute APIs
- [X] Implement "get VM instance" (by ID)
    - [ ] Fill all fields from VM instance
        - [X] Original flavor
        - [ ] More fields
- [X] Implement "list VM instances"
    - [ ] Filter by project ID
    - [ ] Filter by project name
    - [ ] Filter by hypervisor
    - [ ] Filter by availability zone
    - [ ] More filter criteria
- [X] Create connection to Identity APIs
- [X] Implement "get project" (by ID)
    - [X] Fill all fields from project
- [X] Implement "list projects"
    - [X] Add filter criteria
- [X] Check that joins between instance and project (by ID) work
- [ ] ...
