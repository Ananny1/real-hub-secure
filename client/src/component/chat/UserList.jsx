const UserList = ({ users, selectedUser, onSelect }) => {
    return (
        <aside className="user-list">
            <h3>Online Users</h3>
            <ul>
                {users.length === 0 ? (
                    <li>No users online</li>
                ) : (
                    users.map((user, index) => (
                        <li
                            key={user.id}
                            className={`user-list-item ${user.nickname === selectedUser ? 'active' : ''}`}
                            onClick={() => onSelect(user.nickname)}
                        >
                            🟢 {user.nickname}
                        </li>
                    ))
                )}
            </ul>
        </aside>
    );
};

export default UserList;
