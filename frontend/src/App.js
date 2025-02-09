import TableComponent from './components/TableComponent';
import './App.css';
import { useEffect, useState } from 'react';
import axios from 'axios'


function App() {
  const [data, setData] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  const apiUrl = process.env.REACT_APP_API_URL;

  useEffect(() => {
    axios.get(`${apiUrl}/pings`)
      .then(response => {
        setData(response.data);
        setLoading(false);
      })
      .catch(error => {
        console.error(error);
        setError('Failed to fetch data');
        setLoading(false);
      });
  }, []);

  return (
    <div className="App">
      {loading && <div>Loading...</div>}
      {error && <div className="error-message">{error}</div>}
      {!loading && !error && <TableComponent data={data} />}
    </div>
  );
}

export default App;
