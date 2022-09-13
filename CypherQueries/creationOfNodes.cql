create (est1:Student{id:'1'})
create (est2:Student{id:'2'})
create (est3:Student{id:'3'})
create (est4:Student{id:'4'})
create (est5:Student{id:'5'})
create (est6:Student{id:'6'})



create (tech:Teacher{id:'1'})

create (grp1:Group{id:'1'})
create (grp2:Group{id:'2'})

create (sem:Semester{year:'2022-2'})


create(mat:Subject{id:'1'})

create(ingsistemas:Carrer{id:'1'})
create(compscience:Carrer{id:'2'})

create
(est1)-[:Registered]->(grp1),
(est2)-[:Registered]->(grp1),
(est3)-[:Registered]->(grp1),
(est4)-[:Registered]->(grp1),

(est6)-[:Registered]->(grp1)

create
(est2)-[:Registered]->(grp2),
(est3)-[:Registered]->(grp2),
(est4)-[:Registered]->(grp2),
(est5)-[:Registered]->(grp2)

create
(teach)-[:Teaches]->(grp1),
(teach)-[:Teaches]->(grp2)

create
(grp1)-[:Belongs]->(mat),
(grp2)-[:Belongs]->(mat)


create
(grp1)-[:Happends]->(sem),
(grp2)-[:Happends]->(sem)


create 
(est1)-[:Enrolled]->(ingsistemas),
(est1)-[:Enrolled]->(compscience),

(est2)-[:Enrolled]->(ingsistemas),
(est2)-[:Enrolled]->(compscience),

(est3)-[:Enrolled]->(ingsistemas),

(est4)-[:Enrolled]->(ingsistemas),

(est6)-[:Enrolled]->(ingsistemas)
(est6)-[:Enrolled]->(compscience),

//? DELETE ALL NODES
MATCH (n)
DETACH DELETE n