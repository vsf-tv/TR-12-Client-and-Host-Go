import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration, ConfigurationOptions } from '../configuration'
import type { Middleware } from '../middleware';

import { Aes128 } from '../models/Aes128';
import { Aes256 } from '../models/Aes256';
import { Channel } from '../models/Channel';
import { ChannelConfiguration } from '../models/ChannelConfiguration';
import { ChannelState } from '../models/ChannelState';
import { ChannelStatus } from '../models/ChannelStatus';
import { ChannelType } from '../models/ChannelType';
import { ConfigurationData } from '../models/ConfigurationData';
import { ConnectRequestContent } from '../models/ConnectRequestContent';
import { ConnectResponseContent } from '../models/ConnectResponseContent';
import { Connection } from '../models/Connection';
import { DeprovisionResponseContent } from '../models/DeprovisionResponseContent';
import { DeviceConfiguration } from '../models/DeviceConfiguration';
import { DeviceRegistration } from '../models/DeviceRegistration';
import { DeviceStatus } from '../models/DeviceStatus';
import { DisconnectResponseContent } from '../models/DisconnectResponseContent';
import { DtlsFingerprint } from '../models/DtlsFingerprint';
import { DtlsSetupRole } from '../models/DtlsSetupRole';
import { EncryptionAes } from '../models/EncryptionAes';
import { EncryptionAes128 } from '../models/EncryptionAes128';
import { EncryptionAes256 } from '../models/EncryptionAes256';
import { EnumValues } from '../models/EnumValues';
import { ErrorDetails } from '../models/ErrorDetails';
import { GetConfigurationResponseContent } from '../models/GetConfigurationResponseContent';
import { GetConnectionStatusResponseContent } from '../models/GetConnectionStatusResponseContent';
import { Health } from '../models/Health';
import { HealthLevel } from '../models/HealthLevel';
import { IceParameters } from '../models/IceParameters';
import { IceServer } from '../models/IceServer';
import { IdAndValue } from '../models/IdAndValue';
import { Profile } from '../models/Profile';
import { ProfileDefinition } from '../models/ProfileDefinition';
import { RangeValues } from '../models/RangeValues';
import { ReportActualConfigurationRequestContent } from '../models/ReportActualConfigurationRequestContent';
import { ReportActualConfigurationResponseContent } from '../models/ReportActualConfigurationResponseContent';
import { ReportStatusRequestContent } from '../models/ReportStatusRequestContent';
import { ReportStatusResponseContent } from '../models/ReportStatusResponseContent';
import { RistCaller } from '../models/RistCaller';
import { RistCallerTransportProtocol } from '../models/RistCallerTransportProtocol';
import { RistListener } from '../models/RistListener';
import { RistListenerTransportProtocol } from '../models/RistListenerTransportProtocol';
import { RistStreamIdentifier } from '../models/RistStreamIdentifier';
import { Rtp } from '../models/Rtp';
import { RtpFecConfiguration } from '../models/RtpFecConfiguration';
import { RtpFecStreamConfig } from '../models/RtpFecStreamConfig';
import { RtpTransportProtocol } from '../models/RtpTransportProtocol';
import { Setting } from '../models/Setting';
import { SettingProfile } from '../models/SettingProfile';
import { SettingsChoice } from '../models/SettingsChoice';
import { SimpleSettings } from '../models/SimpleSettings';
import { SrtCaller } from '../models/SrtCaller';
import { SrtCallerTransportProtocol } from '../models/SrtCallerTransportProtocol';
import { SrtListener } from '../models/SrtListener';
import { SrtListenerTransportProtocol } from '../models/SrtListenerTransportProtocol';
import { StatusValue } from '../models/StatusValue';
import { StreamId } from '../models/StreamId';
import { SupportedProtocol } from '../models/SupportedProtocol';
import { SynchronizationSource } from '../models/SynchronizationSource';
import { Thumbnail } from '../models/Thumbnail';
import { TransportProtocol } from '../models/TransportProtocol';
import { WebRtc } from '../models/WebRtc';
import { WebRtcFecConfig } from '../models/WebRtcFecConfig';
import { WebRtcFecMechanism } from '../models/WebRtcFecMechanism';
import { WebRtcTransportProtocol } from '../models/WebRtcTransportProtocol';
import { ZixiCaller } from '../models/ZixiCaller';
import { ZixiCallerTransportProtocol } from '../models/ZixiCallerTransportProtocol';
import { ZixiListener } from '../models/ZixiListener';
import { ZixiListenerTransportProtocol } from '../models/ZixiListenerTransportProtocol';

import { ObservableDefaultApi } from "./ObservableAPI";
import { DefaultApiRequestFactory, DefaultApiResponseProcessor} from "../apis/DefaultApi";

export interface DefaultApiConnectRequest {
    /**
     * 
     * @type ConnectRequestContent
     * @memberof DefaultApiconnect
     */
    connectRequestContent: ConnectRequestContent
}

export interface DefaultApiDeprovisionRequest {
    /**
     * 
     * Defaults to: undefined
     * @type string
     * @memberof DefaultApideprovision
     */
    hostId: string
    /**
     * 
     * Defaults to: undefined
     * @type boolean
     * @memberof DefaultApideprovision
     */
    force?: boolean
}

export interface DefaultApiDisconnectRequest {
}

export interface DefaultApiGetConfigurationRequest {
}

export interface DefaultApiGetConnectionStatusRequest {
}

export interface DefaultApiReportActualConfigurationRequest {
    /**
     * 
     * @type ReportActualConfigurationRequestContent
     * @memberof DefaultApireportActualConfiguration
     */
    reportActualConfigurationRequestContent: ReportActualConfigurationRequestContent
}

export interface DefaultApiReportStatusRequest {
    /**
     * 
     * @type ReportStatusRequestContent
     * @memberof DefaultApireportStatus
     */
    reportStatusRequestContent: ReportStatusRequestContent
}

export class ObjectDefaultApi {
    private api: ObservableDefaultApi

    public constructor(configuration: Configuration, requestFactory?: DefaultApiRequestFactory, responseProcessor?: DefaultApiResponseProcessor) {
        this.api = new ObservableDefaultApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param param the request object
     */
    public connectWithHttpInfo(param: DefaultApiConnectRequest, options?: ConfigurationOptions): Promise<HttpInfo<ConnectResponseContent>> {
        return this.api.connectWithHttpInfo(param.connectRequestContent,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public connect(param: DefaultApiConnectRequest, options?: ConfigurationOptions): Promise<ConnectResponseContent> {
        return this.api.connect(param.connectRequestContent,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public deprovisionWithHttpInfo(param: DefaultApiDeprovisionRequest, options?: ConfigurationOptions): Promise<HttpInfo<DeprovisionResponseContent>> {
        return this.api.deprovisionWithHttpInfo(param.hostId, param.force,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public deprovision(param: DefaultApiDeprovisionRequest, options?: ConfigurationOptions): Promise<DeprovisionResponseContent> {
        return this.api.deprovision(param.hostId, param.force,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public disconnectWithHttpInfo(param: DefaultApiDisconnectRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<DisconnectResponseContent>> {
        return this.api.disconnectWithHttpInfo( options).toPromise();
    }

    /**
     * @param param the request object
     */
    public disconnect(param: DefaultApiDisconnectRequest = {}, options?: ConfigurationOptions): Promise<DisconnectResponseContent> {
        return this.api.disconnect( options).toPromise();
    }

    /**
     * @param param the request object
     */
    public getConfigurationWithHttpInfo(param: DefaultApiGetConfigurationRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<GetConfigurationResponseContent>> {
        return this.api.getConfigurationWithHttpInfo( options).toPromise();
    }

    /**
     * @param param the request object
     */
    public getConfiguration(param: DefaultApiGetConfigurationRequest = {}, options?: ConfigurationOptions): Promise<GetConfigurationResponseContent> {
        return this.api.getConfiguration( options).toPromise();
    }

    /**
     * @param param the request object
     */
    public getConnectionStatusWithHttpInfo(param: DefaultApiGetConnectionStatusRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<GetConnectionStatusResponseContent>> {
        return this.api.getConnectionStatusWithHttpInfo( options).toPromise();
    }

    /**
     * @param param the request object
     */
    public getConnectionStatus(param: DefaultApiGetConnectionStatusRequest = {}, options?: ConfigurationOptions): Promise<GetConnectionStatusResponseContent> {
        return this.api.getConnectionStatus( options).toPromise();
    }

    /**
     * @param param the request object
     */
    public reportActualConfigurationWithHttpInfo(param: DefaultApiReportActualConfigurationRequest, options?: ConfigurationOptions): Promise<HttpInfo<ReportActualConfigurationResponseContent>> {
        return this.api.reportActualConfigurationWithHttpInfo(param.reportActualConfigurationRequestContent,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public reportActualConfiguration(param: DefaultApiReportActualConfigurationRequest, options?: ConfigurationOptions): Promise<ReportActualConfigurationResponseContent> {
        return this.api.reportActualConfiguration(param.reportActualConfigurationRequestContent,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public reportStatusWithHttpInfo(param: DefaultApiReportStatusRequest, options?: ConfigurationOptions): Promise<HttpInfo<ReportStatusResponseContent>> {
        return this.api.reportStatusWithHttpInfo(param.reportStatusRequestContent,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public reportStatus(param: DefaultApiReportStatusRequest, options?: ConfigurationOptions): Promise<ReportStatusResponseContent> {
        return this.api.reportStatus(param.reportStatusRequestContent,  options).toPromise();
    }

}
