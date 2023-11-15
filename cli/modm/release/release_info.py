from typing import Any
from msrest.serialization import Model


class PlanInfo(Model):
    _attribute_map = {
        "name": {"key": "name", "type": "str"},
        "publisher": {"key": "publisher", "type": "str"},
        "product": {"key": "product", "type": "str"},
    }

    def __init__(self, name=None, publisher=None, product=None):
        super().__init__()
        self.name = name
        self.publisher = publisher
        self.product = product


class ImageReference(Model):
    _attribute_map = {
        "publisher": {"key": "publisher", "type": "str"},
        "offer": {"key": "offer", "type": "str"},
        "sku": {"key": "sku", "type": "str"},
        "version": {"key": "version", "type": "str"},
    }

    def __init__(self, publisher=None, offer=None, sku=None, version=None, **kwargs: Any):
        super().__init__(**kwargs)
        self.publisher = publisher
        self.offer = offer
        self.sku = sku
        self.version = version


class OfferInfo(Model):
    _attribute_map = {"plan": {"key": "plan", "type": "PlanInfo"}, "image_reference": {"key": "imageReference", "type": "ImageReference"}}

    def __init__(self, plan=None, image_reference=None, **kwargs: Any):
        super().__init__(**kwargs)
        self.plan: PlanInfo = plan
        self.image_reference: ImageReference = image_reference


class ReferenceInfo(Model):
    _attribute_map = {"offer": {"key": "offer", "type": "OfferInfo"}, "vmi": {"key": "vmi", "type": "str"}}

    def __init__(self, offer=None, vmi=None, **kwargs: Any):
        super().__init__(**kwargs)
        self.offer: OfferInfo = offer
        self.vmi = vmi


class ResourcesInfo(Model):
    """This class is used to represent a release's resources. This corresponds to the "resources" property in the index.json at the root of the repo"""

    _attribute_map = {
        "download_url": {"key": "downloadUrl", "type": "str"},
        "filename": {"key": "filename", "type": "str"},
        "sha256_digest": {"key": "sha256Digest", "type": "str"},
    }

    def __init__(self, download_url=None, filename=None, sha256_digest=None):
        super(ResourcesInfo, self).__init__()
        self.download_url = download_url
        self.filename = filename
        self.sha256_digest = sha256_digest


class ReleaseInfo(Model):
    """
    Release information

    Schema:
        conforms to object represented in the releases collection defined in ./schemas/releaseIndex.json
    """

    _attribute_map = {
        "version": {"key": "version", "type": "str"},
        "description": {"key": "description", "type": "str"},
        "reference": {"key": "reference", "type": "ReferenceInfo"},
        "resources": {"key": "resources", "type": "ResourcesInfo"},
    }

    def __init__(self, version=None, description=None, reference=None, resources=None):
        super(ReleaseInfo, self).__init__()
        self.version = version
        self.description = description
        self.reference = reference
        self.resources: ResourcesInfo = resources
