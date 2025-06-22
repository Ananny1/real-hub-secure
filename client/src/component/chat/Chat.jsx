import { useState, useEffect, useRef } from "react";
import UserList from "./UserList";
import ChatBox from "./ChatBox";
import MessageInput from "./MessageInput";
import './chat.css';

export default function Chat() {
    const [users, setUsers] = useState([]);
    const [messages, setMessages] = useState({});
    const [selectedUser, setSelectedUser] = useState(null);
    const wsRef = useRef(null);

    // 1. Fetch user list on mount
    useEffect(() => {
        fetch("http://localhost:8080/chat/users", { credentials: "include" })
            .then(res => res.json())
            .then(data => setUsers(data || []))
            .catch(() => setUsers([]));
    }, []);

    // 2. Select the first user when user list is loaded
    useEffect(() => {
        if (users.length > 0) setSelectedUser(users[0]);
    }, [users]);

    // 3. Fetch chat history whenever selectedUser changes
    useEffect(() => {
        if (!selectedUser) return;
        fetch(`http://localhost:8080/chat/history?with=${selectedUser.id}`, {
            credentials: "include"
        })
            .then(res => res.json())
            .then(history => {
                const safeHistory = Array.isArray(history) ? history : [];
                setMessages(prev => ({
                    ...prev,
                    [selectedUser.id]: safeHistory.map(msg => ({
                        from: msg.sender_id === selectedUser.id
                            ? (users.find(u => u.id === selectedUser.id)?.nickname || "")
                            : "",
                        text: msg.content,
                        imageUrl: msg.image_url
                    }))
                }));
            })
            .catch(() => {
                setMessages(prev => ({ ...prev, [selectedUser.id]: [] }));
            });
    }, [selectedUser, users]);

    // 4. WebSocket setup for real-time messaging
    useEffect(() => {
        wsRef.current = new window.WebSocket("ws://localhost:8080/ws");

        wsRef.current.onopen = () => console.log("WS open!");
        wsRef.current.onclose = () => console.log("WS close!");
        wsRef.current.onerror = err => console.log("WS error", err);

        wsRef.current.onmessage = msg => {
            try {
                const data = JSON.parse(msg.data);
                // Handle online user list update
                if (data.type === "userlist") {
                    setUsers(data.users || []);
                    if (
                        selectedUser &&
                        !data.users.some(u => u.id === selectedUser.id)
                    ) {
                        setSelectedUser(null);
                    }
                }
                // Handle incoming chat message (real-time)
                if (data.type === "chat") {
                    setMessages(prev => ({
                        ...prev,
                        [data.from]: [
                            ...(prev[data.from] || []),
                            {
                                from: data.nickname || "Them",
                                text: data.content,
                                imageUrl: data.image_url
                            }
                        ]
                    }));
                }
            } catch (e) {
                console.log("WS bad data:", msg.data);
            }
        };

        return () => {
            if (wsRef.current) wsRef.current.close();
        };
        // eslint-disable-next-line
    }, [selectedUser]);

    // 5. Send a text message to the selected user
    const handleSendMessage = (text) => {
        if (!selectedUser || !text.trim()) return;
        
        // Add outgoing message locally
        setMessages(prev => ({
            ...prev,
            [selectedUser.id]: [
                ...(prev[selectedUser.id] || []),
                { from: "", text, imageUrl: null }
            ]
        }));
        
        // Send over WebSocket
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({
                type: "chat",
                to: selectedUser.id,
                content: text,
                image_url: ""
            }));
        }
    };

    // 6. Send an image to the selected user
    const handleSendImage = (imageUrl) => {
        if (!selectedUser || !imageUrl) return;
        
        // Add outgoing image message locally
        setMessages(prev => ({
            ...prev,
            [selectedUser.id]: [
                ...(prev[selectedUser.id] || []),
                { from: "", text: "", imageUrl }
            ]
        }));
        
        // Send over WebSocket
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({
                type: "chat",
                to: selectedUser.id,
                content: "",
                image_url: imageUrl
            }));
        }
    };

    return (
        <main className="chat-page">
            <div className="chat-wrapper">
                <UserList
                    users={users}
                    selectedUser={selectedUser?.nickname}
                    onSelect={nickname => {
                        const userObj = users.find(u => u.nickname === nickname);
                        setSelectedUser(userObj);
                    }}
                />
                <div className="chat-main">
                    <ChatBox messages={messages[selectedUser?.id] || []} />
                    {selectedUser && (
                        <MessageInput 
                            onSend={handleSendMessage} 
                            onSendImage={handleSendImage}
                        />
                    )}
                </div>
            </div>
        </main>
    );
}