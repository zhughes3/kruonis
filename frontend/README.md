This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Hey there other person!

Thanks for helping out with the project,I hope you have a blast working on it. :) <br />
Please allow me to tell you a couple of things about the project.

### Project Structure

The project is still pretty small and has a very basic structure, as seen and explained below. 

#### `Project tree`

```bash
├───Assets
│   └───Icons
├───Components
├───Containers
├───Http
├───Interfaces
│   └───Enums
└───Utils
```

#### `Assets`

Contains all images, fonts, icons, etc.

#### `Components`

Every page (`Container`) consists of components. For example: A page contains a form.
The page can be found in the `Containers` folder. The form can be found here, in the `Components` folder.

#### `Containers`

The `Containers` folder holds the root of a page. So when you navigate to the homepage, the code for that page starts in:
`Containers -> Home.tsx`
A synonym for the word Container would be: Page.

#### `Http`

This project uses the fetch api to make http requests. This option was chosen because the API is present natively in all
major browsers, and no additional packages will have to be installed.

##### `Structure`

This folder contains two files: `HttpSetup.ts`, and `Requests.ts`.

1. The `HttpSetup.ts` file contains all basic information we need to create a request (url, port, etc). It also contains 
    basic `GET`, `PUT`, `POST`, `DELETE` functions that all function in the `Requests.ts` file (described below) use.

2.  The `Requests.ts` file contains all requests made to the api. Not every available endpoint on the API has a request yet,
    so you may need to add some.
    
#### `Interfaces`

Because this project uses typescript, we can add interfaces to our objects, props etc. Please add all interfaces to this folder. <br />
<b>Please use Interfaces wherever possible! It makes code much more maintainable and easier to read.</b>

##### `Enums`
    
Likewise, please add all enums to this folder. <br />
<b>Please use Enums wherever possible! It makes code much more maintainable and easier to read.</b>

#### `Utils`

This folder contains functions or variables that are used in multiple files.
For example: password checks are stored in there, since these are checked on creating an account and resetting a password.

## Redux

Redux has not been added to the project yet. So far there has been no reason to. State is management within components,
and never travels more than 1 component down or up.

If the situation arises where data needs to traverse 2 or more components, please add redux to handle those state changes.

## Available Scripts

In the project directory, you can run:

### `yarn start`

Runs the app in the development mode.<br />
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.<br />
You will also see any lint errors in the console.

### `yarn test`

Launches the test runner in the interactive watch mode.<br />
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `yarn build`

Builds the app for production to the `build` folder.<br />
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br />
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `yarn eject`

**Note: this is a one-way operation. Once you `eject`, you can’t go back!**

If you aren’t satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you’re on your own.

You don’t have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn’t feel obligated to use this feature. However we understand that this tool wouldn’t be useful if you couldn’t customize it when you are ready for it.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).
