const toml = require('toml');
const fs = require('fs');
const handlebars = require('handlebars');
const handlebarsIntl = require('handlebars-intl');
handlebarsIntl.registerWith(handlebars);

const defaultFile = "./resume.toml";
const htmlTemplate = "./html.hbs";

function readFile(filename) {
  return new Promise((res, rej) => {
    fs.readFile(filename, 'utf-8', (err, data) => {
      if (err) rej(err);
      res(data);
    });
  });
}

function writeFile(filename, data) {
  return new Promise((res, rej) => {
    fs.writeFile(filename, data, {encoding: 'utf-8'}, (err) => {
      if (err) rej(err);
      return res();
    });
  });
}

async function compileHTML() {
  try {
    const tomlData = await readFile(defaultFile);
    const resumeData = toml.parse(tomlData);

    const htmlTemplateData = await readFile(htmlTemplate);
    const htmlTemplateFn = handlebars.compile(htmlTemplateData);

    return htmlTemplateFn(resumeData);
  } catch (e) {
    console.error(`Unable to open ${defaultFile} for reading: ${e}`);
  }
}

async function main() {
  html = await compileHTML()
  writeFile("resume.html", html);
}

main();


