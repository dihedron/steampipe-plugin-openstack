# steampipe-plugin-openstack

A Steampipe plugin to query OpenStack data as you would a relational database.

At the moment you can run queries such as:

```sql
select 
    port.id, 
    port.device_id, 
    port.device_owner, 
    proj.id, 
    proj.name, 
    vm.id, 
    vm.name 
from 
    openstack_port port 
join 
    openstack_project proj 
on 
    proj.id = port.project_id 
left outer join 
    openstack_instance vm 
on 
    port.device_id = vm.id;
```

Runninga query is as simple as:

```bash
$> steampipe query "select * from openstack_instance where id = 'foo';"

$> steampipe query "select vm.id, vm.name, vm.host_id, vm.flavor_sockets, vm.flavor_disk, prj.name, prj.enabled, prj.id from openstack_instance vm, openstack_project prj where vm.id = '12345678-90ab-cdef-1234-567890abcdef' and vm.project_id = prj.id;"
```

# TODO

This plugin is still in the very early stages.

- [x] Skeleton
- [x] Configuration schema
    - [X] Make OpenStack API micro-versions configurable
- [X] Create connection to OpenStack APIs
- [X] Create connection to Compute APIs
- [X] Implement "get VM instance" (by ID)
    - [ ] Fill all fields from VM instance
        - [X] Original flavor
        - [ ] More fields
- [X] Implement "list VM instances"
    - [X] Filter by project ID
    - [X] Filter by hypervisor
    - [X] Filter by availability zone
    - [ ] Filter by tags
- [X] Create connection to Identity APIs
- [X] Implement "get project" (by ID)
    - [X] Fill all fields from project
- [X] Implement "list projects"
    - [X] Add filter criteria
    - [ ] Filter by tags

- [X] Create connection to Network APIs
- [X] Implement "get port" (by ID)
    - [X] Fill all fields from port
- [X] Implement "list ports"
    - [X] Add filter criteria
    - [ ] Filter by tags    
- [X] Check that joins between instance and project (by ID) work
- [ ] Understand how to expunge embedded entities from instance (e.g . []SecurityGroups)
- [ ] ...
