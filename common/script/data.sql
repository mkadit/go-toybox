-- Sample data insertion
INSERT INTO roles (name, description, is_default) VALUES  ('user', 'Standard user with basic permissions', true);
INSERT INTO roles (name, description, is_default) VALUES  ('admin', 'Administrator with full system access', false);
INSERT INTO roles (name, description, is_default) VALUES  ('moderator', 'Content moderator with limited admin rights', false);
-- INSERT INTO roles (name, description, is_default) VALUES 

-- URL management permissions
INSERT INTO permissions (name, description, category) VALUES  ('create_short_url', 'Ability to create short URLs', 'url_management');
INSERT INTO permissions (name, description, category) VALUES  ('edit_own_url', 'Edit own created URLs', 'url_management');
INSERT INTO permissions (name, description, category) VALUES  ('delete_own_url', 'Delete own created URLs', 'url_management');
INSERT INTO permissions (name, description, category) VALUES  ('view_own_urls', 'View own created URLs', 'url_management');
-- User management permissions
INSERT INTO permissions (name, description, category) VALUES  ('manage_users', 'Ability to manage user accounts', 'user_management');
INSERT INTO permissions (name, description, category) VALUES  ('view_user_list', 'View list of users', 'user_management');
INSERT INTO permissions (name, description, category) VALUES  ('edit_user_roles', 'Change user roles', 'user_management');

-- System permissions
INSERT INTO permissions (name, description, category) VALUES  ('view_system_logs', 'View system logs', 'system');
INSERT INTO permissions (name, description, category) VALUES  ('manage_system_settings', 'Modify system settings', 'system');
-- INSERT INTO permissions (name, description, category) VALUES ;


---- Sample role-permission mappings
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    r.id, 
    p.id 
FROM 
    roles r, 
    permissions p
WHERE 
    (r.name = 'user' AND p.name IN ('create_short_url', 'edit_own_url', 'delete_own_url', 'view_own_urls'))
