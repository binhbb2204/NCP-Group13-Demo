-- Test script to verify friends query
-- Based on your data: id=1, user_id=3, friend_id=2, status=accepted

USE chat_sys;

-- Show all users
SELECT 'All Users:' as info;
SELECT * FROM users;

-- Show all friendships
SELECT 'All Friendships:' as info;
SELECT * FROM friendships;

-- Test query for user_id = 2 (should see user 3 as friend)
SELECT 'Friends of User 2:' as info;
SELECT DISTINCT u.id, u.username,
    (SELECT COUNT(*) FROM messages 
     WHERE sender_id = u.id AND recipient_id = 2 AND is_read = 0) as unread_count
FROM users u
INNER JOIN friendships f ON 
    (f.user_id = 2 AND f.friend_id = u.id) OR 
    (f.friend_id = 2 AND f.user_id = u.id)
WHERE f.status = 'accepted' AND u.id != 2
ORDER BY u.username;

-- Test query for user_id = 3 (should see user 2 as friend)
SELECT 'Friends of User 3:' as info;
SELECT DISTINCT u.id, u.username,
    (SELECT COUNT(*) FROM messages 
     WHERE sender_id = u.id AND recipient_id = 3 AND is_read = 0) as unread_count
FROM users u
INNER JOIN friendships f ON 
    (f.user_id = 3 AND f.friend_id = u.id) OR 
    (f.friend_id = 3 AND f.user_id = u.id)
WHERE f.status = 'accepted' AND u.id != 3
ORDER BY u.username;
