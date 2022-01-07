class CV {
  static example() {
    return CV(
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
      detials: Detials(
        name: "Meneer banaanmans",
        reference: "d1b77757-144a-4e67-871d-3a3c0f910c9b",
        email: "some-email@example.com",
        phoneNr: "06 11111111",
        birthDate: DateTime.now(),
        streetName: "Straatnaam",
        houseNumber: "15",
        houseNumberSuffix: "a",
        zip: "1234AB",
        city: "Groningen",
        country: "Nederland",
        scrapedFromWebsite: "werk.nl",
      ),
      workExpr: [
        WorkExp(
          "Gangster",
          "Gangsters B.V.",
          DateTime.utc(2000, 1, 1),
          DateTime.utc(2010, 1, 1),
          "Doing that gangster shit you know",
        ),
        WorkExp(
          "Bel meneer",
          "De bel mensen",
          DateTime.utc(1995, 2, 1),
          DateTime.utc(1998, 5, 1),
          "Je weet wel die hinderlijke mensen die je opbellen met onzin waar je helemaal geen belang bij hebt",
        ),
      ],
      education: [
        Education(
          "HBO",
          "Hogeschool",
          DateTime.utc(2000, 1, 1),
          DateTime.utc(2010, 1, 1),
        ),
        Education(
          "MBO",
          "Some school",
          DateTime.utc(1995, 2, 1),
          DateTime.utc(1998, 5, 1),
        ),
      ],
      courses: [
        Education(
          "beeing smart",
          "Hogeschool",
          DateTime.utc(2010, 1, 1),
          DateTime.utc(2011, 1, 1),
        ),
      ],
      driversLicenses: ["B", "C", "D"],
      languageSkills: [
        LanguageSkill(
          "Nederlands",
          LanguageSkillLevel.excellent,
          LanguageSkillLevel.good,
        ),
        LanguageSkill(
          "Duits",
          LanguageSkillLevel.good,
          LanguageSkillLevel.reasonable,
        ),
      ],
    );
  }

  CV({
    required this.createdAt,
    required this.updatedAt,
    required this.detials,
    required this.workExpr,
    required this.education,
    required this.courses,
    required this.languageSkills,
    required this.driversLicenses,
  });

  final DateTime createdAt;
  final DateTime updatedAt;
  final Detials detials;
  final List<WorkExp> workExpr;
  final List<Education> education;
  final List<Education> courses;
  final List<LanguageSkill> languageSkills;
  final List<String> driversLicenses;
}

class Detials {
  Detials({
    required this.name,
    required this.reference,
    required this.email,
    required this.phoneNr,
    required this.birthDate,
    required this.streetName,
    required this.houseNumber,
    required this.houseNumberSuffix,
    required this.zip,
    required this.city,
    required this.country,
    required this.scrapedFromWebsite,
  });

  final String name;
  final String reference;
  final String email;
  final String phoneNr;
  final DateTime birthDate;
  final String streetName;
  final String houseNumber;
  final String houseNumberSuffix;
  final String zip;
  final String city;
  final String country;
  final String scrapedFromWebsite;

  bool get hasAddress =>
      streetName != null && houseNumber != null && city != null;
}

enum LanguageSkillLevel {
  unknown,
  reasonable,
  good,
  excellent,
}

Map<LanguageSkillLevel, String> humanLanguageSkillLevel = {
  LanguageSkillLevel.unknown: "Onbekend",
  LanguageSkillLevel.reasonable: "Redelijk",
  LanguageSkillLevel.good: "Goed",
  LanguageSkillLevel.excellent: "Uitstekend",
};

Map<LanguageSkillLevel, int> languageSkillLevelNr = {
  LanguageSkillLevel.unknown: 0,
  LanguageSkillLevel.reasonable: 1,
  LanguageSkillLevel.good: 2,
  LanguageSkillLevel.excellent: 3,
};

const int maxLanguageLevelNr = 3;

class LanguageSkill {
  const LanguageSkill(this.name, this.reading, this.writing);

  final String name;
  final LanguageSkillLevel reading;
  final LanguageSkillLevel writing;
}

class Education {
  const Education(this.name, this.org, this.from, this.to);

  final String name;
  final String org;
  final DateTime from;
  final DateTime to;
}

class WorkExp {
  const WorkExp(
    this.name,
    this.company,
    this.from,
    this.to,
    this.description,
  );

  final String name;
  final String company;
  final DateTime from;
  final DateTime to;
  final String description;
}
