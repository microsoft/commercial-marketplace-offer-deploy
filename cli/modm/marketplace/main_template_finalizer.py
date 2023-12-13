import copy
from modm.release.release_info import ReferenceInfo
from .main_template import MainTemplate


class MainTemplateFinalizer:
    def __init__(self, main_template: MainTemplate) -> None:
        self.main_template = main_template

    def finalize(self, **kwargs):
        """
        Updates the (installer's) main template with the parameters from the app's main template.
        This results in a flow of: createUiDefinition.json/parameters/outputs --> mainTemplate.json/parameters

        Explanation:
            This allows the parameters to be passed to the mainTemplate.json/variables/userData
            so MODM can bootstrap the application with it's parameters when it performs the deployment
        """
        main_template = copy.deepcopy(self.main_template)
        main_template.set_parameters(kwargs.get("template_parameters", []))
        main_template.user_data.set_installer_package_hash(kwargs.get("installer_package").hash)

        release_info: ReferenceInfo = kwargs.get("release_info")
        use_vmi_reference = kwargs.get("use_vmi_reference", False)
        vmi_reference_id = kwargs.get("vmi_reference_id", None)

        if use_vmi_reference:
            if vmi_reference_id is not None:
                main_template.vmi_reference_id = vmi_reference_id
            else:
                main_template.vmi_reference_id = release_info.vmi
        else:
            if release_info is None:
                raise ValueError("reference_info is required when not passing a VMI reference")
            main_template.vm_offer = release_info.reference.offer.serialize()

        return main_template
