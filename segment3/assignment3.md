## Assignment 3 - Pet shelter

You are building a management system for a pet shelter. The shelter keeps track of the animals in its care and helps match them with new owners.

Your program should have two structs: a `Pet` and a `Shelter`.

A `Pet` has a name, a species, an age, and a status indicating whether they have been adopted.

A `Shelter` has a name and keeps a list of all the pets currently in its care.

To get started, populate your shelter with at least 5 creative and original pets of your own invention.

The user interacts with the shelter through a menu:
1. List available pets
2. Adopt a pet
3. Add a new pet
4. Find the oldest available pet
5. Quit

Behaviour:

- Listing available pets shows only pets that have not yet been adopted, with their full profile — name, species, age, and status.
- Adopting a pet takes a name as input, marks that pet as adopted, and prints a farewell message. If no pet with that name exists, or they are already adopted, print an appropriate message.
- Adding a new pet asks the user for a name, species, and age.
- Finding the oldest available pet prints the profile of the non-adopted pet with the highest age. If all pets have been adopted, print an appropriate message.

Methods to consider:

- A method on `Pet` to print its profile
- A method on `Pet` to mark it as adopted
- A method on `Shelter` to add a pet
- A method on `Shelter` to list available pets
- A method on `Shelter` to adopt a pet by name
- A method on `Shelter` to find the oldest available pet

Think carefully about which methods need a pointer receiver and which do not.
