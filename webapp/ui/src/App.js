import {
  createBrowserRouter,
  RouterProvider,
  Route,
  Outlet,
} from "react-router-dom";
import Register from "./views/Register";
import Login from "./views/Login";
import Index from "./views/Index";
import AllPackages from "./views/AllPackages";
import SinglePackage from "./views/SinglePackage";
import Navigator from "./components/Navigator";
import "./style.scss"


const PageLayout = () => {
  return (
    <>
      <Navigator/>
      <Outlet/>
    </>
  );
} 

const browserRouter = createBrowserRouter([
  {
    path: "/",
    element: <PageLayout/>,
    children: [
      {
        path: "/",
        element: <Index/>,
      },
      {
        path: "/package/all",
        element: <AllPackages/>,
      },
      {
        path: "/package/:id",
        element: <SinglePackage/>,
      },

    ]
  },
  {
    path: "/login",
    element: <Login/>,
  },
  {
    path: "/register",
    element: <Register/>,
  },
]);

function App() {
  return (
    <div className = "app">
      <div className = "container">
        <RouterProvider router={browserRouter}/>
      </div>
    </div>
  );
}

export default App;