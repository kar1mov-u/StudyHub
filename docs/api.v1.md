## Status
v1-Draft

Modules - represnts a university module, timeless entity:
    GET /modules   - retrieves all modules, can add filetring based on Departmets
    POST /modules    - creates a module, can be only called by admin
    GET /modules/<module_id>    - retrieves the whole info about the module and current moduleRun and list of teaching weeks
    DELETE /modules<module_id>    - deletes the module and all data related to it

Module-Run  -- represents the module cohort,cause 1 module can be taught in multiple semesters,for each cohort will have its own moduleRun, identified by id and semester taught
    GET /modules/<module_id>/runs           for a given module , retrieve all module runs. Only 1 module run can be active, meaining is being taught in current semester
    POST /modules/<module_id>/runs           creates a new run for a module.  (Should be called only by the admin in the start of each semester)
    GET /module_runs/<module_run_id>        returns the specified module run, and it will consist of N Weeks
    DELETE /module_runs/<module_run_id>     deletes the moduleRun

Weeks    - represents academic weeks

    GET /weeks/<week_id> - returns the specified week (for now doesnt decide what will be included in the week)


## GET /modules
Purpose: Home page — list all available modules

```[
  {
    "id": "uuid",
    "code": "CS204",
    "name": "Data Structures & Algorithms",
    "department": "Computer Science"
  },
  {
    "id": "uuid",
    "code": "CS201",
    "name": "DataBase Systems",
    "department": "Computer Science"
  }

]```





