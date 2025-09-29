# Grader-Butler
Convenience tool for downloading student lab submissions and getting them on COELinux unzipped and named. The release is compiled for COELinux and intended to be used there.

build: env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/graderb

## Installation (On COELinux)

Download the latest release directly:

```bash
wget https://github.com/<your-username>/<your-repo>/releases/latest/download/<your-binary-or-archive>
```
## Setup

This program requires a config file to run. All values are strings

```json
{
    "auth_token": "...",
    "course_id": "...",
    "assignment_id": "...",
    "lab_title": "...",
    "students": [
        "John Doe",
        "Jane Doe",
        ...
    ]
}
```
### Here are the MANDATORY steps to fill out the config
1. Generate an authentication token. To do this, go to Account -> Settings on Carmen. In the "Approved Integrations" section, click "+ New Access Token". Configure desired settings and paste it in the "auth_token" field (inside the "")
```json
"auth_token": "<your-generated-token>"
```
2. Obtain the ID of the course (Systems 1). To do this, navigate to Courses -> (desired section of Systems 1), and obtain the ID at the end of the URL. For example, if the URL is "https://osu.instructure.com/courses/123456", you want 
```json
"course_id: "123456"
```
3. Obtain the ID of the assignment. To do this, navigate to Courses -> (desired section of Systems 1) -> Assignments -> (desired lab to grade), and obtain the ID at the end of the URL. For example, if the URL is "https://osu.instructure.com/courses/123456/assignments/7654321", you want 
```json 
"assignment_id": "7654321" 
```
4. The lab title field is optional and used for compiling which isn't fully implemented yet. Leave this as "" for now.
5. Build out the array of student names as they appear in the gradebook or Grader Distribution Sheet (Kirby's course). For example, if I wanted to grade John Foo, John Bar, and John Baz, I'd have 
```json
"students": [
    "John Foo",
    "John Bar",
    "John Baz"
]
```

That's it! you should be ready to download student submissions now by simply invoking
```bash
./<name>
```
