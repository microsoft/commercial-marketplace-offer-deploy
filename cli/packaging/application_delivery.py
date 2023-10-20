

from enum import Enum

class ApplicationDelivery(str, Enum):
    """How the application is delivered to the end user."""
    MARKETPLACE = "marketplace"
    SERVICECATALOG = "servicecatalog"
