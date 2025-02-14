# Frontend technical specs

## Pages

### Dashboard /dashboard

#### Purpose
The purpose of this page is to provide a summary of learning and activity as the default page when a user vists the web-app.

#### Components
This page contains the following components:

- Last Study Session
    shows last activity used
    shows when last activity was
    shows how many words were correct vs wrong from last activity
    has a link to the word group
- Study Progress
    - total words study eg. 3/432
        - across all study sessions show the total words studied out of all possible words in our database
    - display a mastery progress 
- Quick Stats
    - success rate eg. 80%
    - total study sessions eg. 5
    - total active groups eg. 3
    - study streak eg. 3 days
    - average study duration eg. 10 min
    - average success rate eg. 80%
- Start Studying button
    - goes to study activittes page

We'll need the following API endpoints to pawer this page:


#### Needed API endpoints
- GET /dashboard/last_study_session
- GET /dashboard/study_progress
- GET /dashboard/quick_stats
- POST /dashboard/start_studying


### Study Activites Page

#### Purpose
The purpose of this page is to show a collection of study activities with a thumbnail and a name that the user can launch

#### Components
- Study Activity Card
    - Show a thumbnail of the study activity
    - Show the name of the study activity
    - Show the description of the study activity
    - Launch button
    - The view page to view more information about pas study sessions for this study activity

#### Needed API endpoints
- GET /study_activities

## Study Activity Show `/study_activities/:id`

#### Purpose
This page is used to show a detailed view of a specific study activity and its past study sessions.

#### Components
- Show the name of the study activity
- Show the description of the study activity
- Show the thumbnail of the study activity
- Launch button
- Study Activities Paginated List
    - id
    - activity name
    - group name
    - start time
    - end time (inferred by the last word_review_item submitted)
    - number of review items

    #### Needed API endpoints
    - GET /api/study_activities/:id
    - GET /api/study_activities/:id/study_sessions

### Study Activites Launch Page `/study_activities/:id/launch`

#### Purpose
This page is used to launch a specific study activity.

#### Components
- Name of study activity
- Launch form
    - select field for group
    - launch now button

#### Behavior
After the form is submitted a new tab opens with the study activity based on its URL provided in the database.
Also after the form is submitted the page will rediret to the study session show page.

#### Needed API endpoints
- POST /api/study_activities/:id/launch


### Words `/words`

#### Purpose
The purpose of this page is to show all words in the database.

#### Components
- Paginated Word List
    - Columns:
        - Chinese character
        - Pinyin
        - English definition
        - Correct Count
        - Incorrect Count
    - Pagination with 100 items per page
    - Clicking the Chinese character will take us to the word show page

#### Needed API endpoints
- GET /api/words

### Word Show `/words/:id`

#### Purpose
This page is used to show information about a specific word.

#### Components
- Chinese character
- Pinyin
- English definition
- Study Statistics
    - Correct Count
    - Incorrect Count
- Word Groups
    - show a series of pill tags
    - when group name is clicked it will take us to the group show page

#### Needed API endpoints
- GET /api/words/:id

### Word Groups Index `/groups`

#### Purpose
This page is used to show all word groups in the database.

#### Components
- Paginated Word Group List
    - Columns:
        - Group Name
        - Description
        - Word Count
    - Pagination with 100 items per page
    - Clicking the Group Name will take us to the group show page

#### Needed API endpoints
- GET /api/groups

### Word Group Show `/groups/:id`

#### Purpose
This page is used to show information about a specific word group.

#### Components
- Group Name
- Group Statistics
    - Word Count
- Words in Group (Paginated list of words)
    - Should use the same component as the words index page
- Study Sessions (Paginated list of study sessions)
    - Should use the same component as the study sessions index page

#### Needed API endpoints
- GET /api/groups/:id (the name and groups stats)
- GET /api/groups/:id/words
- GET /api/groups/:id/study_sessions


### Study Sessions `/study_sessions`

#### Purpose
This page is used to show all study sessions in the database.

#### Components
- Paginated Study Session List
    - Columns:
        - id
        - activity name
        - group name
        - start time
        - end time
        - number of review items
    - Clicking the study session id will take us to the study session show page

#### Needed API endpoints
- GET /api/study_sessions


### Study Session Show `/study_sessions/:id`

#### Purpose
This page is used to show information about a specific study session.

#### Components
- Study Session Details
    - activity name
    - group name
    - start time
    - end time
    - number of review items
- Word Review Items (Paginated list of review items)
    - Should use the same component as the words index page

#### Needed API endpoints
- GET /api/study_sessions/:id
- GET /api/study_sessions/:id/words

### Settings `/settings`

#### Purpose
This page is used to make configurations to the study portal.

#### Components
- Theme Selection eg. Light, Dark, System Default
- Reset History Button
    - This will delete ll study sessions and word review items
- Full Reset Button
    - This will drop all tables and re-create them from scratch

#### Needed API endpoints
- POST /api/settings/theme
- POST /api/settings/reset_history
- POST /api/settings/full_reset








