MATCH (n)
DETACH DELETE n;

CREATE (:Person {id: -12, userId: -12, name: 'Lena', surname: 'Lenić', picture: '', bio: '', quote: '', email: 'autor2@gmail.com'});
CREATE (:Person {id: -13, userId: -13, name: 'Sara', surname: 'Sarić', picture: '', bio: '', quote: '', email: 'autor3@gmail.com'});
CREATE (:Person {id: -22, userId: -22, name: 'Mika', surname: 'Mikić', picture: '', bio: '', quote: '', email: 'turista2@gmail.com'});
CREATE (:Person {id: -23, userId: -23, name: 'Steva', surname: 'Stević', picture: '', bio: '', quote: '', email: 'turista3@gmail.com'});
CREATE (:Person {id: -21, userId: -21, name: 'Pera', surname: 'Perić', picture: '', bio: '', quote: '', email: 'turista1@gmail.com'});