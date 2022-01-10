class CV {
  static example() {
    return CV(
      referenceNumber: "d1b77757-144a-4e67-871d-3a3c0f910c9b",
      createdAt: DateTime.now(),
      lastChanged: DateTime.now(),
      personalDetails: PersonalDetails(
        firstName: "Meneer",
        surName: "banaanmans",
        email: "some-email@example.com",
        phoneNumber: "06 11111111",
        dob: DateTime.now(),
        zip: "1234AB",
        streetName: "Straatnaam",
        houseNumber: "15",
        houseNumberSuffix: "a",
        city: "Groningen",
        country: "Nederland",
      ),
      workExperiences: [
        WorkExperience(
          profession: "Gangster",
          employer: "Gangsters B.V.",
          startDate: DateTime.utc(2000, 1, 1),
          endDate: DateTime.utc(2010, 1, 1),
          description: "Doing that gangster shit you know",
        ),
        WorkExperience(
          profession: "Bel meneer",
          employer: "De bel mensen",
          startDate: DateTime.utc(1995, 2, 1),
          endDate: DateTime.utc(1998, 5, 1),
          description:
              "Je weet wel die hinderlijke mensen die je opbellen met onzin waar je helemaal geen belang bij hebt",
        ),
      ],
      educations: [
        Education(
          "HBO",
          institute: "Hogeschool",
          startDate: DateTime.utc(2000, 1, 1),
          endDate: DateTime.utc(2010, 1, 1),
        ),
        Education(
          "MBO",
          institute: "Some school",
          startDate: DateTime.utc(1995, 2, 1),
          endDate: DateTime.utc(1998, 5, 1),
        ),
      ],
      courses: [
        Course(
          "beeing smart",
          institute: "Hogeschool",
          startDate: DateTime.utc(2010, 1, 1),
          endDate: DateTime.utc(2011, 1, 1),
        ),
        Course(
          "TV host voor dummies",
          institute: "Hogeschool",
          startDate: DateTime.utc(2010, 1, 1),
          // endDate: DateTime.utc(2011, 1, 1),
        ),
      ],
      driversLicenses: ["B", "C", "D"],
      languages: [
        Language(
          "Nederlands",
          LanguageLevel.good,
          LanguageLevel.excellent,
        ),
        Language(
          "Duits",
          LanguageLevel.reasonable,
          LanguageLevel.good,
        ),
      ],
    );
  }

  CV.fromJson(Map<String, dynamic> json) {
    referenceNumber = json['referenceNumber'];
    createdAt = jsonParseDate(json['createdAt']);
    lastChanged = jsonParseDate(json['lastChanged']);
    educations = json['educations']
        ?.map((entry) => Education.fromJson(entry))
        ?.toList()
        ?.cast<Education>();
    courses = json['courses']
        ?.map((entry) => Course.fromJson(entry))
        ?.toList()
        ?.cast<Course>();
    workExperiences = json['workExperiences']
        ?.map((entry) => WorkExperience.fromJson(entry))
        ?.toList()
        ?.cast<WorkExperience>();
    preferredJobs = json['preferredJobs']?.cast<String>();
    languages = json['languages']
        ?.map((entry) => Language.fromJson(entry))
        ?.toList()
        ?.cast<Language>();
    personalDetails = PersonalDetails.fromJson(json['personalDetails']);
    driversLicenses = json['driversLicenses']?.cast<String>();
  }

  CV({
    required this.referenceNumber,
    this.createdAt,
    this.lastChanged,
    this.educations,
    this.courses,
    this.workExperiences,
    this.preferredJobs,
    this.languages,
    required this.personalDetails,
    this.driversLicenses,
  });

  late final String referenceNumber;
  late final DateTime? createdAt;
  late final DateTime? lastChanged;
  late final List<Education>? educations;
  late final List<Course>? courses;
  late final List<WorkExperience>? workExperiences;
  late final List<String>? preferredJobs;
  late final List<Language>? languages;
  late final PersonalDetails personalDetails;
  late final List<String>? driversLicenses;
}

class Education {
  Education.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    description = json['description'];
    institute = json['institute'];
    isCompleted = json['isCompleted'];
    hasDiploma = json['hasDiploma'];
    startDate = jsonParseDate(json['startDate']);
    endDate = jsonParseDate(json['endDate']);
  }

  Education(
    this.name, {
    this.description,
    this.institute,
    this.isCompleted,
    this.hasDiploma,
    this.startDate,
    this.endDate,
  });

  late final String name;
  late final String? description;
  late final String? institute;
  late final bool? isCompleted;
  late final bool? hasDiploma;
  late final DateTime? startDate;
  late final DateTime? endDate;
}

class Course {
  Course.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    institute = json['institute'];
    startDate = jsonParseDate(json['startDate']);
    endDate = jsonParseDate(json['endDate']);
    isCompleted = json['isCompleted'];
    description = json['description'];
  }

  Course(
    this.name, {
    this.institute,
    this.startDate,
    this.endDate,
    this.isCompleted,
    this.description,
  });

  late final String name;
  late final String? institute;
  late final DateTime? startDate;
  late final DateTime? endDate;
  late final bool? isCompleted;
  late final String? description;
}

class WorkExperience {
  WorkExperience.fromJson(Map<String, dynamic> json) {
    profession = json['profession'];
    description = json['description'];
    startDate = jsonParseDate(json['startDate']);
    endDate = jsonParseDate(json['endDate']);
    stillEmployed = json['stillEmployed'];
    employer = json['employer'];
    weeklyHoursWorked = json['weeklyHoursWorked'];
  }

  WorkExperience({
    this.profession,
    this.description,
    this.startDate,
    this.endDate,
    this.stillEmployed,
    this.employer,
    this.weeklyHoursWorked,
  });

  late final String? profession;
  late final String? description;
  late final DateTime? startDate;
  late final DateTime? endDate;
  late final bool? stillEmployed;
  late final String? employer;
  late final int? weeklyHoursWorked;
}

enum LanguageLevel {
  unknown,
  reasonable,
  good,
  excellent,
}

Map<LanguageLevel, String> humanLanguageLevel = {
  LanguageLevel.unknown: "Onbekend",
  LanguageLevel.reasonable: "Redelijk",
  LanguageLevel.good: "Goed",
  LanguageLevel.excellent: "Uitstekend",
};

Map<LanguageLevel, int> languageLevelToNr = {
  LanguageLevel.unknown: 0,
  LanguageLevel.reasonable: 1,
  LanguageLevel.good: 2,
  LanguageLevel.excellent: 3,
};

Map<int, LanguageLevel> nrToLanguageLevel = {
  0: LanguageLevel.unknown,
  1: LanguageLevel.reasonable,
  2: LanguageLevel.good,
  3: LanguageLevel.excellent,
};

const int maxLanguageLevelNr = 3;

class Language {
  Language.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    levelSpoken =
        nrToLanguageLevel[json['levelSpoken']] ?? LanguageLevel.unknown;
    levelWritten =
        nrToLanguageLevel[json['levelWritten']] ?? LanguageLevel.unknown;
  }

  Language(this.name, this.levelSpoken, this.levelWritten);

  late final String name;
  late final LanguageLevel levelSpoken;
  late final LanguageLevel levelWritten;
}

class Competence {
  Competence(this.name, {this.description});

  String name;
  String? description;
}

class Interest {
  Interest(this.name, {this.description});

  String name;
  String? description;
}

class PersonalDetails {
  PersonalDetails({
    this.initials,
    this.firstName,
    this.surNamePrefix,
    this.surName,
    this.dob,
    this.gender,
    this.streetName,
    this.houseNumber,
    this.houseNumberSuffix,
    this.zip,
    this.city,
    this.country,
    this.phoneNumber,
    this.email,
  });

  PersonalDetails.fromJson(Map<String, dynamic> json) {
    initials = json['initials'];
    firstName = json['firstName'];
    surNamePrefix = json['surNamePrefix'];
    surName = json['surName'];
    dob = jsonParseDate(json['dob']);
    gender = json['gender'];
    streetName = json['streetName'];
    houseNumber = json['houseNumber'];
    houseNumberSuffix = json['houseNumberSuffix'];
    zip = json['zip'];
    city = json['city'];
    country = json['country'];
    phoneNumber = json['phoneNumber'];
    email = json['email'];
  }

  late final String? initials;
  late final String? firstName;
  late final String? surNamePrefix;
  late final String? surName;
  late final DateTime? dob;
  late final String? gender;
  late final String? streetName;
  late final String? houseNumber;
  late final String? houseNumberSuffix;
  late final String? zip;
  late final String? city;
  late final String? country;
  late final String? phoneNumber;
  late final String? email;

  bool get hasAddress =>
      streetName != null && houseNumber != null && city != null;

  String get fullName {
    String resp = this.firstName ?? '';

    if (this.surName != null) {
      if (resp.isNotEmpty) resp += " ";

      resp += (this.surNamePrefix != null ? this.surNamePrefix! + " " : '') +
          this.surName!;
    }

    return resp;
  }
}

DateTime? jsonParseDate(dynamic input) {
  if (input == null) return null;
  return DateTime.parse(input);
}
