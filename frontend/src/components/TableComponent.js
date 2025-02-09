import React from 'react';
import { Table } from 'antd';
import { format } from 'date-fns';

const columns = [
    {
        title: 'IP адрес',
        dataIndex: 'ip_address',
        key: 'ip',
    },
    {
        title: 'Время пинга',
        dataIndex: 'ping_time',
        key: 'pingTime',
        render: (text) => formatDate(text),
    },
    {
        title: 'Статус',
        dataIndex: 'is_success',
        key: 'is_success',
        render: (is_success) => (is_success ? '✅ Успешный' : '❌ Неуспешный'),
    },
];

const formatDate = (dateString) => {
    const date = new Date(dateString);
    return format(date, 'yyyy-MM-dd HH:mm:ss');
};

const TableComponent = ({ data }) => (
    <div>
        <h1>Статус контейнеров</h1>
        <Table columns={columns} dataSource={data} />
    </div>
);

export default TableComponent;
