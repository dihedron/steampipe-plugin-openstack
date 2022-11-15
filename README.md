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

- *General understanding*
    - [ ] Understand how to expunge into own table embedded entities (e.g. []SecurityGroups from port)
    - [ ] Understand how to work with JSONb data
    - [ ] Understand how to use different operators in list filter (e.g. "description match regexp")
- *Implementation* 
    - [x] Skeleton
    - [x] Configuration schema
        - [X] Make OpenStack API micro-versions configurable
    - [X] Create connection
        - [X] OpenStack APIs
        - [X] Identity APIs
        - [X] Compute APIs
        - [X] Network APIs
        - [X] Block Storage APIs
    - [X] Virtual machines instances
        - [X] Get
        - [X] List
            - [X] Filter
        - [X] Embed original flavor
        - [ ] *TODO*
            - [ ] Add more fields
            - [ ] Embed image info
            - [ ] Manage tags
    - [X] Network ports
        - [X] Get
        - [X] List
            - [X] Filter
        - [ ] *TODO*
            - [ ] Manage tags
    - [X] Projects
        - [X] Get
        - [X] List
            - [X] Filter
        - [ ] *TODO*
            - [ ] Manage tags
    - [X] Block storage volumes
        - [X] Get
        - [X] List
            - [X] Filter
        - [ ] *TODO*
            - [ ] Manage metadata
    - [X] Check that joins between entities work
