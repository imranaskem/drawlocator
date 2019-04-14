let res = [
    db.createCollection("people"),
    db.people.insert({
        id: "1",
        firstname: "Kent",
        lastname: "Valentine",
        placeofwork: "Holiday"
    }),
    db.people.insert({
        id: "2",
        firstname: "Dean",
        lastname: "Faulkner",
        placeofwork: "Holiday"
    }),
    db.people.insert({
        id: "3",
        firstname: "Matt",
        lastname: "Tobin",
        placeofwork: "Holiday"
    }),
    db.people.insert({
        id: "4",
        firstname: "Imran",
        lastname: "Askem",
        placeofwork: "Holiday"
    }),
    db.people.insert({
        id: "5",
        firstname: "Sian",
        lastname: "Barlow",
        placeofwork: "Holiday"
    }),
    db.people.insert({
        id: "6",
        firstname: "John",
        lastname: "Sheard",
        placeofwork: "Holiday"
    })
]

printjson(res)