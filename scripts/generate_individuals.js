import { faker } from "@faker-js/faker";
import fs from "fs";
import path from "path";

function generateIndividual() {
  return {
    id: faker.string.uuid(),
    first_name: faker.person.firstName(),
    last_name: faker.person.lastName(),
    birth_date: faker.date.birthdate().toISOString().split("T")[0],
    death_date: faker.datatype.boolean()
      ? faker.date.past({ years: 10 }).toISOString().split("T")[0]
      : "",
    gender: faker.person.sex(),
    events: faker.helpers.arrayElements(
      ["birth", "marriage", "death", "graduation", "job"],
      faker.number.int({ min: 0, max: 3 })
    ),
  };
}

function generateIndividuals(count) {
  const arr = [];
  for (let i = 0; i < count; i++) arr.push(generateIndividual());
  return arr;
}

const filePath = path.join("data", "individuals.json");
const data = generateIndividuals(200);

fs.writeFileSync(filePath, JSON.stringify(data, null, 2), "utf-8");
console.log("OK ->", filePath);
