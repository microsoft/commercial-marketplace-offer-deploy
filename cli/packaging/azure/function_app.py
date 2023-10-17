import uuid

def create_function_app_name(prefix: str):
    """
    Creates a function app name based on the prefix.
    """
    return f"{prefix}{uuid.uuid4().hex[:12].lower()}"